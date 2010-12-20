#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>

void		launch_job()
{
	int	p[2];
	char	*argv[2][3] = {
		{"ls", NULL},
		{"cat", "-e", NULL}
	};

	pipe(p);
	if (fork() == 0) {
		dup2(p[1], 1);
		close(p[1]);
		execvp(argv[0][0], argv[0]);
		exit(1);
	} else {
		close(p[1]);
		dup2(p[0], 0);
		close(p[0]);
		execvp(argv[1][0], argv[1]);
		exit(1);
	}
}

int main() {
	launch_job();
	return (0);
}
