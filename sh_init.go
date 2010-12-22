package main

import (
	"os"
	"strconv"
)

/**
 * @brief Copy current env in local one
 */
func (self *Gosh) loadEnv() {
	self.env = make([]string, len(os.Environ()))
	copy(self.env, os.Environ())
}

/**
 * @brief Update SHLVL env variable
 */
func (self *Gosh) updateShlvl() {
	var (
		lvl    int
		lvlStr string
		err    os.Error
	)

	if lvlStr, _, err = self.getEnv("SHLVL"); err != nil {
		lvlStr = "0"
	}
	if lvl, err = strconv.Atoi(lvlStr); err != nil {
		lvl = 0
	}
	self.setEnv("SHLVL", strconv.Itoa(lvl+1))
}

