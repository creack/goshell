package main

import (
	"fmt"
	"strings"
	"syscall"
)

/**
 * @brief Check the command and add path if needed
 *
 * @param argv List of arguments of the command
 *
 */
func (self *Gosh) cmdCheckPath(argv []string) {

	/// First we check if the command is a builting
	if fctBuiltin, check := self.builtins[argv[0]]; check == true {
		fctBuiltin(self, argv)
		return
	}

	/// If here, it means, it is not a builtin, if it does not start with
	/// '/', './' or '../' try to concat it with path
	directPrefix := []string{"/", "./", "../"}
	for _, elem := range directPrefix {
		if strings.HasPrefix(argv[0], elem) {
			self.exec(argv[0], argv)
			return
		}
	}

	/// If it is a "regular" command, try the path
	path, _, err := self.getEnv("PATH")
	if err != nil {
		path = DEFAULT_PATH
	}
	/// @todo Use os.Stat/Permission instead of syscall.Access
	pathTab := strings.Split(path, ":", -1)
	for _, elem := range pathTab {
		if syscall.Access(elem+"/"+argv[0], 0) == 0 {
			self.exec(elem+"/"+argv[0], argv)
			return
		}
	}
	fmt.Printf("gosh: command not found: %s\n", argv[0])
}

/**
 * @brief Parse the line from stdin
 *
 * @param line Line read from stdin
 *
 */
func (self *Gosh) parse(line string) {
	argv := strings.Fields(line)
	self.cmdCheckPath(argv)
}


