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
 * @brief [Builtin] Display env
 */
func env(sh *Gosh) {
	for _, line := range sh.env {
		fmt.Printf("%s\n", line)
	}
}

/**
 * @brief [Builtin] Exit shell
 */
func exit(sh *Gosh) {
	fmt.Printf("Exit\n")
	os.Exit(0)
}

/**
 * @brief [Builtin] Display current pwd
 */
func getPwd(sh *Gosh) {
	var (
		pwd string
		err os.Error
	)
	if pwd, err = sh.getEnv("PWD"); err != nil {
		if pwd, err = os.Getwd(); err != nil {
			fmt.Printf("Error: can't retrieve pwd\n")
		}
	}
	fmt.Printf("%s\n", pwd)
}

/**
 * @brief Put the builtins functions in object map
 *
 * @todo Use method pointer instead of function pointer
 *
 * @return Map with builtin => function
 */
func defineBuiltins() map[string]func(*Gosh) {
	b := make(map[string]func(*Gosh))

	b["env"] = env
	b["exit"] = exit
	b["pwd"] = getPwd
	return b
}
