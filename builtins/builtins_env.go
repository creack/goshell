package main

import (
	"fmt"
	"strings"
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
		if line != "" {
			fmt.Printf("%s\n", line)
		}
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
	if len(argv) == 1 {
		env(sh, argv)
		return
	}
	index := strings.Index(argv[1], "=")
	if index < 0 {
		fmt.Printf("Error, expected \"export key=value\"\n")
		return
	}
	sh.setEnv(argv[1][:index], argv[1][index+1:])
}

/**
 * @brief [Builtin] Remove specifed var from env
 *
 * @param sh Instance of the shell
 * @param argv Argument list of the command
 *
 */
func unsetEnv(sh *Gosh, argv []string) {
	if len(argv) == 1 {
		return
	}
	if len(argv) == 2 && argv[1] == "*" {
		sh.env = make([]string, 1)
	}
}

/**
 * @brief Display env array size
 */
func envSize(sh *Gosh, argv []string) {
	used := 0
	for i := 0; i < len(sh.env); i++ {
		if sh.env[i] != "" {
			used++
		}
	}
	fmt.Printf("Used : %d/%d\n", used, len(sh.env))
}

