/**
 * @file main.go
 * @brief Main of GoSh
 * @author Guillaume J. CHARMES
 * @version 0.01
 * @date 2010-12-19
 */
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
)

/// default path in case there is no $PATH in env
const (
	DEFAULT_PATH = "/bin:/usr/bin:/usr/local/bin:/opt/bin:/opt/local/bin"
)

type Gosh struct {
	env      []string
	pRead    chan string
	builtins map[string]builtinFunc
}

/**
 * @brief Instanciate the shell
 *
 * @note Should not be called more than once
 *
 * @return New instance of the shell
 */
func NewGosh() *Gosh {
	return &Gosh{
		pRead:    make(chan string),
		builtins: defineBuiltins(),
	}
}

/*
t_sigs        gl_errors[] =
     {
       {SIGSEGV, "Segmentation Fault"},
       {SIGBUS, "Bus Error"},
       {SIGINT, "User Interupted"},
       {SIGQUIT, "Quit"},
       {SIGILL, "Illegal Instruction"},
       {SIGABRT, "Abort"},
       {SIGKILL, "Kill"},
       {SIGTRAP, "Trap"},
       {SIGTERM, "Term"},
       {SIGFPE, "Floating exception"},
       {SIGSYS, "Unknown system call"},
       {SIGPIPE, "Broken pipe"},
       {0, 0}
     };
*/

/**
 * @brief Execute argv[0]
 *
 * @todo Handle errors, Pass correct flags to os.Wait instead of 0
 * @todo Think about pipeline/jobcontrol
 * @todo Check if os.ForkExec is pertinent
 *
 * @param cmd Command to execute with full path
 * @param argv Array of args, argv[0] is the command to execute
 *
 */
func (self *Gosh) exec(cmd string, argv []string) {
	if pid, err := fork(); err != nil || pid < 0 {
		log.Exitf("Error fork: %s\n", err)
	} else if pid == 0 {
		if err = os.Exec(cmd, argv, self.env); err != nil {
			log.Exitf("Error exec: %s\n", err)
		}
	} else {
		var (
			waitStatus *os.Waitmsg
			err        os.Error
		)
		if waitStatus, err = os.Wait(pid, 0); err != nil {
			log.Printf("Error wait: %s\n", err)
			return
		}
		if waitStatus.WaitStatus != 0 {
			fmt.Fprintf(os.Stderr, "%s\n", signal.UnixSignal(waitStatus.WaitStatus.Signal()))
		}
	}
}

/**
 * @brief Laumch the shell
 *
 * @note Use goroutine in order to get stdin, should not.
 *
 */
func (self *Gosh) Start() {
	buf := make([]byte, 1024)
	self.loadEnv()
	self.updateShlvl()
	for {
		print("$>")
		if n, err := os.Stdin.Read(buf); err == os.EOF {
			fmt.Printf("Exit\n")
			break
		} else if err != nil {
			log.Exitf("Error: %s\n", err)
		} else {
			line := string(buf[:n-1])
			if line = strings.TrimSpace(line); line != "" {
				self.parse(line)
			}
		}
	}
}

/**
 * @brief Main
 */
func main() {
	sh := NewGosh()
	sh.Start()
}
