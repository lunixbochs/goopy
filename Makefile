include $(GOROOT)/src/Make.inc

TARG=vm
GOFILES=\
		vm.go \
		builtins.go \
		object.go \
		frame.go \
		types/int.go \
		types/bool.go \
		types/string.go \
		types/type.go \
		types/none.go


include $(GOROOT)/src/Make.cmd
