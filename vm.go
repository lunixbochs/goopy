package main

import (
	"reflect"
)

type Machine struct {
	Globals		map[string]*Object
	//Builtins	map[string]*Object
	Builtins 	*Builtins
	Frames		[]*Frame // keeps track of the call stack
	CurFrame 	int // starts at -1, ends at -1
}

func NewMachine() Machine {
	m := Machine{make(map[string]*Object), &Builtins{}, make([]*Frame, 0), -1}
	return m
}

func (VM *Machine) Run(f *Frame) {
	f.VM = VM
	VM.Frames = append(VM.Frames, f)
	VM.CurFrame++
	cf := VM.Frames[VM.CurFrame]
	for VM.CurFrame >= 0 {
		cf = VM.Frames[VM.CurFrame]
		if cf.Cur < len(cf.Code) && cf.Ret == nil{
			cf.Step()
		} else if cf.Ret != nil {
			VM.CurFrame-- // cut back to our old frame, the function caller
			caller := *VM.Frames[VM.CurFrame]
			caller.Local[caller.RetLocal] = cf.Ret
		} else { // we must be done with this frame, and should regress
			VM.CurFrame--
		}
	}
}

func (VM *Machine) Call(f *Object, args chan *Object) {
	switch f.Value.(type) {
		case NativeFunction: // we're dealing with a builtin wrapper
			f.Value.(NativeFunction)(VM, args)
		case Func: // this is a "real" function written in the interpreted language
			VM.Frames = append(VM.Frames, f.Value.(Func).Frame)
			VM.CurFrame++
			cf := *VM.Frames[VM.CurFrame]
			for i := 0; len(args) != 0; i++ {
				cf.Local[i] = <-args
			}
	}
}

func (VM *Machine) Lookup(name string) (o *Object) {
	if o = VM.Globals[name]; o != nil {
		return o
	}
	v := reflect.ValueOf(VM.Builtins)
	fn := v.MethodByName("B_" + name)
	if fn.IsValid() {
		return fn.Call([]reflect.Value{})[0].Interface().(*Object)
	} else {
		panic("tried to lookup nonexistent name:  "+name)
	}
	return NewObject(nil)
}

func main() {
	 //builtin function call test case
	/*consts := []*Object{&Object{INT, 1}, &Object{INT, 2}}
	names := make([]string, 1)
	names[0] = "print"
	code := []Bytecode{
		{CONST, 1, 0, 0},
		{CONST, 2, 1, 0},
		{ADD, 3, 1, 2},
		{ARGPUSH, 0, 3, LOCAL},
		{GET, 4, 0, 0},
		{CALL, 0, 4, 0},
		{ARGPUSH, 0, 0, LOCAL},
		{CALL, 0, 4, 0}}*/
	 // int, raw_input, IF test case
	consts := []*Object{NewObject(3), NewObject("less than 3"), NewObject("equal to 3"), NewObject("greater than 3")}
	names := make([]string, 3)
	names[0] = "raw_input"
	names[1] = "print"
	names[2] = "int"
	code := []Bytecode{
		{GET, 1, 0, 0},
		{CALL, 2, 1, 0},
		{GET, 1, 0, 2},
		{ARGPUSH, 0, 2, LOCAL},
		{CALL, 2, 1, 0},
		{CONST, 3, 0, 0},
		{LT, 4, 2, 3},
		{IF, 4, 0, 0},
		{JUMP, 0, 4, 0},
		{ARGPUSH, 0, 1, CONSTANT},
		{GET, 5, 0, 1},
		{CALL, 0, 5, 0},
		{JUMP, 0, 10, 0},
		{EQ, 4, 2, 3},
		{IF, 4, 0, 0},
		{JUMP, 0, 4, 0},
		{ARGPUSH, 0, 2, CONSTANT},
		{GET, 5, 0, 1},
		{CALL, 0, 5, 0},
		{JUMP, 0, 3, 0},
		{ARGPUSH, 0, 3, CONSTANT},
		{GET, 5, 0, 1},
		{CALL, 0, 5, 0},
		{NOP, 0, 0, 0}}
	// while loop test
	/*consts := []*Object{NewObject(10), NewObject(1)}
	names := make([]string, 1)
	names[0] = "print"
	code := []Bytecode{
		{CONST, 1, 0, 0},
		{CONST, 2, 0 , 1},
		{GET, 3, 0, 0},
		{IF, 1, 0, 0},
		{JUMP, 0, 4, 0},
		{SUB, 1, 1, 2},
		{ARGPUSH, 0, 1, LOCAL},
		{CALL, 0, 3, 0},
		{JUMP, 0, 6, 1}}*/
	VM := NewMachine()

	f := NewFrame(consts, names, code)
	VM.Run(f)
}
