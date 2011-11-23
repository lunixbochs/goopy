include $(GOROOT)/src/Make.inc

TARG=vm
GOFILES=\
		op.go \
		vm.go \
		builtins.go \
		object.go \
		frame.go

include $(GOROOT)/src/Make.cmd
