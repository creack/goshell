package main

import (
	"container/list"
	//"fmt"
	"os"
	"strings"
	"syscall"
	//"./signalC/_obj/signalC"
)

/**
 * @brief Check if the given command exists and is executable
 *
 * @todo Use os.Stat/Permission insteas of syscall.Access
 * @todo Check if executable
 *
 * @param cmd Command to test
 *
 * @return Error if not exists, nil if exists
 */
func (self *Gosh) cmdCheckAccess(cmd string) os.Error {
	if syscall.Access(cmd, 0) != 0 {
		return os.NewError("gosh: command not found: ")
	}
	return nil
}

/**
 * @brief Check the command and add path if needed
 *
 * @param cmd Command to check
 *
 * @todo Check if the file exists when begin with / ./ ../
 * @return string with correct path, bool to check if builtin or not and error if any
 */
func (self *Gosh) cmdCheckPath(cmd string) (string, bool, os.Error) {
	var err os.Error

	/// First we check if the command is a builting
	if _, check := self.builtins[cmd]; check == true {
		return cmd, true, nil
	}

	/// If here, it means, it is not a builtin, if it does not start with
	/// '/', './' or '../' try to concat it with path
	directPrefix := []string{"/", "./", "../"}
	for _, elem := range directPrefix {
		if strings.HasPrefix(cmd, elem) {
			return cmd, false, nil
		}
	}

	/// If it is a "regular" command, try the path
	path, _, err := self.getEnv("PATH")
	if err != nil {
		path = DEFAULT_PATH
	}
	pathTab := strings.Split(path, ":", -1)
	for _, elem := range pathTab {
		if err = self.cmdCheckAccess(elem + "/" + cmd); err == nil {
			return elem + "/" + cmd, false, nil
		}
	}
	return "", false, os.NewError(err.String() + cmd)
}

/**
 * @brief Parse the line from stdin
 *
 * Parsing strategy : Split on separtors (; && ||) and then use regexp.
 * It is not really efficient but I think it is going to work.
 *
 * @note This "parsing" is extremly ugly, I know.
 *
 * @todo Handle && ||
 *
 * @param line Line read from stdin
 *
 * @return Job list ready to be executed and error if any
 */
func (self *Gosh) parse(line string) (*jobList, os.Error) {
	var (
		jobs      *jobList
		isBuiltin bool
		err       os.Error
	)

	jobs = &jobList{list.New()}
	colonJobs := strings.Split(line, ";", -1)
	for _, elem := range colonJobs {
		if elem = strings.TrimSpace(elem); elem != "" {
			j := jobs.PushBack(NewJob(elem))
			processPiped := strings.Split(elem, "|", -1)
			for _, pCmd := range processPiped {
				if pCmd = strings.TrimSpace(pCmd); pCmd != "" {
					argv := strings.Fields(pCmd)
					if argv[0], isBuiltin, err = self.cmdCheckPath(argv[0]); err != nil {
						return nil, err
					}
					j.Value.(*job).process.PushBack(NewProcess(argv, isBuiltin))
				}
			}
		}
	}
	return jobs, nil
}
