/**
 * @file sh_env.go
 * @brief Environment relating internal functions
 * @author Guillaume J. CHARMES
 * @version 0.01
 * @date 2010-12-21
 */
package main

import (
	"strings"
	"os"
)

/**
 * @brief Set env variable
 *
 * @param key Name of the variable
 * @param value Value of the variable
 *
 * @todo use Gosh.getEnv in order to check if var exists (see getEnv todo)
 */
func (self *Gosh) setEnv(key, value string) {
	/// We check if the key exists
	for i := 0; i < len(self.env); i++ {
		index := strings.Index(self.env[i], "=")
		if index < 0 {
			continue
		}
		if self.env[i][:index] == key {
			self.env[i] = key+"="+value
			return
		}
	}
	/// If here, it mean the key is new, so we insert it
	var i int
	for i = 0; i < len(self.env); i++ {
		if self.env[i] == "" {
			break
		}
	}
	if i == len(self.env) {
		newEnv := make([]string, len(self.env)*2)
		copy(newEnv, self.env)
		self.env = newEnv
	}
	self.env[i] = key + "=" + value
}

/**
 * @brief Get the env value of the givven key
 *
 * @param key Key (env var) which we want to get
 *
 * @todo make it return the index
 *
 * @return value of the variable and error if any
 */
func (self *Gosh) getEnv(key string) (string, os.Error) {
	for _, line := range self.env {
		index := strings.Index(line, "=")
		if index < 0 {
			continue
		}
		if line[:index] == key {
			return line[index+1:], nil
		}
	}
	return "", os.NewError("Error: $" + key + " is not defined")
}
