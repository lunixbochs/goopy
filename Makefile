include $(GOROOT)/src/Make.inc

TARG=vm
GOFILES=\
		op.go \
		vm.go

include $(GOROOT)/src/Make.cmd
