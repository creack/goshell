/**
* @file defineBuiltins.go
* @brief builtins gosh
* @author Guillaume J. CHARMES
* @version 0.01
* @date 2010-12-19
 */
package main

import (
	"fmt"
	"os"
)

/**
 * @brief [Builtin] Exit shell
 *
 * @param sh Instance of the shell
 * @param argv Argument list of the command
 *
 */
func exit(sh *Gosh, argv []string) {
	fmt.Printf("Exit\n")
	os.Exit(0)
}

/**
 * @brief [Builtin] Display current pwd
 *
 * @param sh Instance of the shell
 * @param argv Argument list of the command
 *
 */
func getPwd(sh *Gosh, argv []string) {
	var (
		pwd string
		err os.Error
	)
	if pwd, _, err = sh.getEnv("PWD"); err != nil {
		if pwd, err = os.Getwd(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: can't retrieve pwd\n")
		}
	}
	fmt.Printf("%s\n", pwd)
}

type builtinFunc func(*Gosh, []string)

/**
 * @brief Put the builtins functions in object map
 *
 * @todo Use method pointer instead of function pointer
 *
 * @return Map with builtin => function
 */
func defineBuiltins() map[string]builtinFunc {
	b := make(map[string]builtinFunc)

	b["env"] = env
	b["getenv"] = getEnv
	b["export"] = setEnv
	b["setenv"] = setEnv
	b["unset"] = unsetEnv
	b["unsetenv"] = unsetEnv
	b["envsize"] = envSize

	b["cd"] = chdir
	b["chdir"] = chdir

	b["exit"] = exit
	b["pwd"] = getPwd
	return b
}
