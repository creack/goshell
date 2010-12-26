package main

import (
	"container/list"
	"fmt"
	"os"
	"strings"
	"syscall"
	"./signalC/_obj/signal"
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
 * @brief Launch job
 *
 * @todo Handle stderr
 * @todo Put this function and all job related in own file
 * @todo Cut this function in subfunctions
 *
 * @param sh Shell instance
 *
 */
func (j *job) start(sh *Gosh) {
	var (
		pid                       int
		infile, outfile           *os.File
		infileOldFd, outfileOldFd int
		readPipe, writePipe       *os.File
		fds                       []*os.File
		err                       os.Error
	)
	fds = make([]*os.File, 3)
	infile = j.stdin
	for e := j.process.Front(); e != nil; e = e.Next() {
		p := e.Value.(*process)

		/// Set up pipes if necessary
		if e.Next() != nil {
			if readPipe, writePipe, err = os.Pipe(); err != nil {
				fmt.Fprintf(os.Stderr, "Error pipe: %s\n", err)
				return
			}
			outfile = writePipe
		} else {
			outfile = j.stdout
		}
		if p.isBuiltin == true {
			var errno int

			if infile != j.stdin {
				if infileOldFd, errno = syscall.Dup(infile.Fd()); errno != 0 {
					fmt.Fprintf(os.Stderr, "Error dup: %d\n", errno)
					return
				}
				if errno = syscall.Dup2(infile.Fd(), j.stdin.Fd()); errno != 0 {
					fmt.Fprintf(os.Stderr, "Error dup2: %d\n", errno)
					return
				}
			}
			if outfile != j.stdout {
				if outfileOldFd, errno = syscall.Dup(j.stdout.Fd()); errno != 0 {
					fmt.Fprintf(os.Stderr, "Error dup: %d\n", errno)
					return
				}
				if errno = syscall.Dup2(outfile.Fd(), j.stdout.Fd()); errno != 0 {
					fmt.Fprintf(os.Stderr, "Error dup2: %d\n", errno)
					return
				}
			}
			sh.builtins[p.argv[0]](sh, p.argv)
			if infile != j.stdin {
				if errno = syscall.Dup2(infileOldFd, j.stdin.Fd()); errno != 0 {
					fmt.Fprintf(os.Stderr, "Error dup2 back: %d\n", errno)
					return
				}
			}
			if outfile != j.stdout {
				if errno = syscall.Dup2(outfileOldFd, j.stdout.Fd()); errno != 0 {
					fmt.Fprintf(os.Stderr, "Error dup2 back: %d\n", errno)
					return
				}
			}
		} else {
			fds[0] = infile
			fds[1] = outfile
			fds[2] = os.Stderr
			signal.RestoreAll()
			if pid, err = os.ForkExec(p.argv[0], p.argv, sh.env, "", fds); err != nil {
				fmt.Fprintf(os.Stderr, "Error fork: %s\n", err)
				return
			}
			if j.pgid == 0 {
				j.pgid = pid
			}
			syscall.Setpgid(pid, j.pgid)
			p.pid = pid
			signal.IgnoreAll()
		}
		if infile != j.stdin {
			infile.Close()
		}
		if outfile != j.stdout {
			if p.isBuiltin == true {
				syscall.Dup2(os.Stdout.Fd(), outfileOldFd)
			}
			outfile.Close()
		}
		infile = readPipe
	}
	j.wait(sh)
}

/**
 * @brief Check all process from all job and mark status
 *
 * @param w WaitMsg return by os.Wait
 *
 * @todo Better error handling, making special error types, etc
 *
 * @return Error if any
 */
func (sh *Gosh) markProcessStatus(w *os.Waitmsg) os.Error {
	if w.Pid <= 0 {
		return os.NewError("This shoudl never append")
	}
	for j := sh.jobList.Front(); j != nil; j = j.Next() {
			for p := j.Value.(*job).process.Front(); p != nil; p = p.Next() {
				if p.Value.(*process).pid == w.Pid {
					if w.Stopped() == true {
						p.Value.(*process).stopped = true
					} else {
						p.Value.(*process).completed = true
						if w.Signaled() == true {
							fmt.Fprintf(os.Stderr, "%d: Terminated by signal %d.\n", w.Pid, w.Signal())
						}
					}
					return nil
				}
			}
	}
	return os.NewError(fmt.Sprintf("No child process %d.\n", w.Pid))
}

/**
 * @brief Check the status of each process of the job. Block until then
 *
 * @param sh Instance of shell
 *
 * @todo Find a proper and protable way to wait on any childs
 * @todo Kill all process in case of wait error
 * @todo Handle errors properly
 */
func (j *job) wait(sh *Gosh) {
	for {
		w, err := os.Wait(-1, os.WUNTRACED)
		if err != nil {
			//fmt.Fprintf(os.Stderr, "Error wait: %s\n", err)
			return
		}
		if sh.markProcessStatus(w); err != nil {
			break
		}
		if j.isStopped() == true {
			break
		}
		if j.isCompleted() == true {
			break
		}
	}
}

/**
 * @brief Check each process if the job is stopped
 *
 * If at least one process is not flagged as completed and stopped,
 * it means that it is still alive.
 *
 * @return true if job is stopped, false otherwise
 */
func (j *job) isStopped() bool {
	for p := j.process.Front(); p != nil; p = p.Next() {
		if p.Value.(*process).completed == false {
			return false
		}
		if p.Value.(*process).stopped == false {
			return false
		}
	}
	return true
}

/**
 * @brief Check each process if the job is completed
 *
 * If at least one process is not flagged as completed, it means
 * the job is not completed/
 *
 * @return true if job is completed, false otherwise
 */
func (j *job) isCompleted() bool {
	for p := j.process.Front(); p != nil; p = p.Next() {
		if p.Value.(*process).completed == false {
			return false
		}
	}
	return true
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
