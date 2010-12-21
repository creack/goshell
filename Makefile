include ${GOROOT}/src/Make.inc

TARG	=	gosh
GOFILES	=	main.go			\
		fork.go			\

GOFILES	+=	builtins/defineBuiltins.go	\
		builtins/builtins_env.go	\

include ${GOROOT}/src/Make.cmd
