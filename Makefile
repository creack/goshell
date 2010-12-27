include ${GOROOT}/src/Make.inc

TARG	=	gosh
GOFILES	=	main.go			\
		sh_init.go		\
		sh_jobs.go		\
		sh_processLine.go	\
		sh_env.go		\
		fork.go			\

GOFILES	+=	builtins/defineBuiltins.go	\
		builtins/builtins_env.go	\
		builtins/builtins_cd.go		\

include ${GOROOT}/src/Make.cmd
