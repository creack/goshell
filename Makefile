include ${GOROOT}/src/Make.inc

TARG	=	gosh
GOFILES	=	main.go			\
		fork.go			\
		defineBuiltins.go	\

include ${GOROOT}/src/Make.cmd
