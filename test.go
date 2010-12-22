package main

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

