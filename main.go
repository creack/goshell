/**
 * @file main.go
 * @brief Main of GoSh
 * @author Guillaume J. CHARMES
 * @version 0.01
 * @date 2010-12-19
 * @todo Make the makefile compile the signalC dep lib
 */
package main

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"./signalC/_obj/signal"
)

/// default path in case there is no $PATH in env
const (
	DEFAULT_PATH = "/bin:/usr/bin:/usr/local/bin:/opt/bin:/opt/local/bin"
)
/** @todo Use map instead of list for procList and Joblib */
type Gosh struct {
	env      []string
	builtins map[string]builtinFunc
	jobList  *jobList
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
		builtins: defineBuiltins(),
		jobList: &jobList{list.New()},
	}
}

type process struct {
	argv               []string /**< For exec */
	pid                int      /**< process ID */
	completed, stopped bool     /**< true if process has completed/stopped */
	status             int      /**< reported status value */
	isBuiltin          bool     /**< If is builtin, no fork */
}
type processList struct {
	*list.List
}

func NewProcess(argv []string, isBuiltin bool) *process {
	p := &process{
		argv:      argv,
		isBuiltin: isBuiltin,
	}
	return p
}

type job struct {
	commandLine           string       /**< command line, used for messages */
	process               *processList /**< list of processes in this job */
	pgid                  int          /**< process group ID */
	notified              bool         /**< true if user told about stopped job */
	stdin, stdout, stderr *os.File     /**< standard i/o channels */
	//tmodes              termios      /**< saved terminal modes */
}
type jobList struct {
	*list.List
}

func NewJob(line string) *job {
	j := &job{
		commandLine: line,
		process:     &processList{list.New()},
		stdin:       os.Stdin,
		stdout:      os.Stdout,
		stderr:      os.Stderr,
	}
	return j
}

/**
 * @brief Execute argv[0]
 *
 * @todo Handle errors, Pass correct flags to os.Wait instead of 0
 * @todo Think about pipeline/jobcontrol
 * @todo Check if os.ForkExec is pertinent
 * @todo Put back the string instead of signal number
 *
 * @param cmd Command to execute with full path
 * @param argv Array of args, argv[0] is the command to execute
 *
 */
func (self *Gosh) exec(cmd string, argv []string) {
	fds := make([]*os.File, 3)
	fds[0] = os.Stdin
	fds[1] = os.Stdout
	fds[2] = os.Stderr
	signal.RestoreAll()
	pid, _ := os.ForkExec(cmd, argv, self.env, "", fds)
	signal.IgnoreAll()

	os.Wait(pid, 0)
	return
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
			fmt.Fprintf(os.Stderr, "%v\n", waitStatus.WaitStatus.Signal())
		}
		_ = waitStatus
	}
}

/**
 * @brief Launch the job
 *
 * @param JobList ready to execute
 *
 */
func (self *Gosh) launchJobs(jobs *jobList) {
	for j := jobs.Front(); j != nil; j = j.Next() {
		j.Value.(*job).start(self)
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
			if jobs, err := self.parse(string(buf[:n-1])); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				continue
			} else {
				self.launchJobs(jobs)
			}
		}
	}
}

/**
 * @brief Main
 */
func main() {
	sh := NewGosh()
	signal.IgnoreAll()
	sh.Start()
}
