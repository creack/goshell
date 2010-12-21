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
	"strings"
	"syscall"
)

/// default path in case there is no $PATH in env
const (
	DEFAULT_PATH = "/bin:/usr/bin:/usr/local/bin:/opt/bin:/opt/local/bin"
)

type Gosh struct {
	env      []string
	pRead    chan string
	builtins map[string]func(*Gosh)
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

/**
 * @brief Read from stdin and send the line to the main loop
 */
func (self *Gosh) Reader() {
	buf := make([]byte, 1024)
	for {
		if n, err := os.Stdin.Read(buf); err != nil {
			log.Exitf("EOF\n")
		} else {
			self.pRead <- string(buf[:n-1])
		}
	}
}

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
		os.Wait(pid, 0)
	}
}

/**
 * @brief Get the env value of the givven key
 *
 * @param key Key (env var) which we want to get
 *
 * @return value of the variable
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
	return "", os.NewError("[env] " + key + " is not defined")
}

/**
 * @brief Check the command and add path if needed
 *
 * @param argv List of arguments of the command
 *
 */
func (self *Gosh) cmdCheckPath(argv []string) {

	/// First we check if the command is a builting
	if fctBuiltin, check := self.builtins[argv[0]]; check == true {
		fctBuiltin(self)
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
	path, err := self.getEnv("PATH")
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

/**
 * @brief Copy current env in local one
 */
func (self *Gosh) loadEnv() {
	osEnv := os.Environ()
	self.env = make([]string, len(osEnv))
	for i, envLine := range osEnv {
		self.env[i] = envLine
	}
}

/**
 * @brief Laumch the shell
 *
 * @note Use goroutine in order to get stdin, should not.
 *
 */
func (self *Gosh) Start() {
	self.loadEnv()
	go self.Reader()
	for {
		print("$>")
		select {
		case line := <-self.pRead:
			if line = strings.TrimSpace(line); line != "" {
				self.parse(line)
			}
		}
	}
}

func toto(recLevel, pid int) {
	var (
		readPipe, writePipe *os.File
		err                 os.Error
	)

	if recLevel >= 5 {
		//if err = os.Exec("/bin/cat", []string{"cat", "-e"}, nil); err != nil {
		//log.Exitf("error exec: %s\n", err)
		//}
		os.Wait(pid, os.WNOHANG|os.WSTOPPED)
		return
	}

	if readPipe, writePipe, err = os.Pipe(); err != nil {
		log.Exitf("Error pipe: %s\n", err)
	}

	if pid, err := fork(); err != nil {
		log.Exitf("Error fork: %s\n", err)
	} else if pid == 0 {
		readPipe.Close()
		syscall.Dup2(writePipe.Fd(), os.Stdout.Fd())
		writePipe.Close()
		if err = os.Exec("/bin/cat", []string{"cat", "-e"}, nil); err != nil {
			log.Exitf("error exec: %s\n", err)
		}
	} else {
		writePipe.Close()
		syscall.Dup2(readPipe.Fd(), os.Stdin.Fd())
		readPipe.Close()
		toto(recLevel+1, pid)
	}
}

func rec_pipe() {
	var (
		readPipe, writePipe *os.File
		err                 os.Error
	)

	if readPipe, writePipe, err = os.Pipe(); err != nil {
		log.Exitf("Error pipe : %s\n", err)
	}
	if pid, err := fork(); err != nil {
		log.Exitf("Error fork: %s\n", err)
	} else if pid == 0 {
		//os.Stdin.Close()
		//os.Stdout = writePipe
		//syscall.Dup2(writePipe.Fd(), os.Stdout.Fd())
		readPipe.Close()
		syscall.Dup2(writePipe.Fd(), os.Stdout.Fd())
		writePipe.Close()
		if err = os.Exec("/bin/ls", []string{"ls"}, nil); err != nil {
			log.Exitf("Error exec son: %s\n", err)
		}
	} else {
		writePipe.Close()
		syscall.Dup2(readPipe.Fd(), os.Stdin.Fd())
		readPipe.Close()
		toto(0, pid)
		//if err = os.Exec("/bin/cat", []string{"cat", "-e"}, nil); err != nil {
		//log.Exitf("error exec: %s\n", err)
		//}
	}
}

/**
 * @brief Main
 */
func main() {
	sh := NewGosh()
	sh.Start()
	<-make(chan int)
	fmt.Printf("Hello World!\n")
}
