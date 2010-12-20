package main

import (
	"os"
	"strconv"
	"syscall"
)

/**
 * @brief Wrap syscall Fork
 *
 * @todo Wrap errors
 *
 * @return Error if any and pid of the son in father proc or 0 in son
 */
func fork() (pid int, err os.Error) {
	darwin := syscall.OS == "darwin"
	r1, r2, err1 := syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if err1 != 0 {
		return 0, os.NewError("Error nº: "+strconv.Itoa(int(err1)))
	}
	// Handle exception for darwin
	if darwin && r2 == 1 {
		r1 = 0;
	}
	return int(r1), nil
}
