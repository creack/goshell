package main

import (
	"fmt"
	"os"
)

/**
 * @brief [Builtin] Implementation of chdir command
 *
 * @param sh Instance of the shell
 * @param argv Argument list of the command
 *
 * @todo handle pwd without os.Getwd
 */
func chdir(sh *Gosh, argv []string) {
	var (
		pwd, oldpwd string
		dest        string
		err         os.Error
	)

	if oldpwd, _, err = sh.getEnv("PWD"); err != nil {
		if oldpwd, err = os.Getwd(); err != nil {
			oldpwd = ""
		}
	}
	if len(argv) == 1 {
		if dest, _, err = sh.getEnv("HOME"); err != nil {
			return
		}
	} else if argv[1] == "-" {
		if dest, _, err = sh.getEnv("OLDPWD"); err != nil {
			return
		}
		fmt.Printf("%s\n", dest)
	} else {
		dest = argv[1]
	}
	if err := os.Chdir(dest); err != nil {
		fmt.Printf("Can't change directory: %s\n", err)
		return
	}
	if pwd, err = os.Getwd(); err != nil {
		pwd = ""
	}
	sh.setEnv("PWD", pwd)
	sh.setEnv("OLDPWD", oldpwd)
}
