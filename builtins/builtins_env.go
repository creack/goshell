package main

import (
	"fmt"
)

/**
 * @brief [Builtin] Display env
 *
 * @param sh Instance of the shell
 * @param argv Argument list of the command
 *
 * @todo Handle env options
 */
func env(sh *Gosh, argv []string) {
	for _, line := range sh.env {
		fmt.Printf("%s\n", line)
	}
}

/**
 * @brief [Builtin] Retrieve specified env variable
 *
 * @param sh Instance of the shell
 * @param argv Argument list of the command
 *
 */
func getEnv(sh *Gosh, argv []string) {
	if len(argv) == 1 {
		env(sh, argv)
		return
	}
	for _, elem := range argv[1:] {
		value, err := sh.getEnv(elem)
		if err != nil {
			value = err.String()
		}
		fmt.Printf("%s\n", value)
	}
}

/**
 * @brief [Builtin] Set env variable
 *
 * @param sh Instance of the shell
 * @param argv Argument list of the command
 *
 */
func setEnv(sh *Gosh, argv []string) {
}

/**
 * @brief [Builtin] Remove specifed var from env
 *
 * @param sh Instance of the shell
 * @param argv Argument list of the command
 *
 */
func unsetEnv(sh *Gosh, argv []string) {
}

