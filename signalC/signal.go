package signalC

// #include <signal.h>
/*
void IgnoreAll() {
	signal(SIGINT,  SIG_IGN);
	signal(SIGQUIT, SIG_IGN);
	signal(SIGTSTP, SIG_IGN);
	signal(SIGTTIN, SIG_IGN);
	signal(SIGTTOU, SIG_IGN);
	signal(SIGCHLD, SIG_IGN);
}

void RestoreAll() {
      signal(SIGINT,  SIG_DFL);
      signal(SIGQUIT, SIG_DFL);
      signal(SIGTSTP, SIG_DFL);
      signal(SIGTTIN, SIG_DFL);
      signal(SIGTTOU, SIG_DFL);
      signal(SIGCHLD, SIG_DFL);
}
*/
import "C"

func IgnoreAll() {
	C.IgnoreAll()
}

func RestoreAll() {
	C.RestoreAll()
}
