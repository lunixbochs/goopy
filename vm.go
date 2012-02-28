package main

import (
	// "fmt"
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
	m := Machine{make(map[string]*Object), &Builtin, make([]*Frame, 0), -1}
	Builtin.B_None = m.NewNone()
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

func (VM *Machine) Call(f *Object, args []*Object) {
	switch f.Value.(type) {
		case NativeFunction: // we're dealing with a builtin wrapper
			f.Value.(NativeFunction)(VM, args)
		case Func: // this is a "real" function written in the interpreted language
			VM.Frames = append(VM.Frames, f.Value.(Func).Frame)
			VM.CurFrame++
			cf := *VM.Frames[VM.CurFrame]
			for i := 0; len(args) != 0; i++ {
				cf.Local[i] = args[i]
			}
		default:
			fn := f.GetAttribute(VM, "__call__")
			if fn != nil {
				VM.Call(fn, args)
			} else {
				panic("tried to call a nil value")
			}
	}
}

func (VM *Machine) RCall(f *Object, args []*Object) *Object{
	ff := FakeFrame(VM)
	VM.Frames = append(VM.Frames, ff)
	VM.CurFrame++
	VM.Call(f, args)
	VM.CurFrame--
	VM.Frames = VM.Frames[:VM.CurFrame+1]
	ret := ff.Local[ff.RetLocal]
	if ret != nil {
		return ret
	}
	return nil//VM.Lookup("None")
}

func (VM *Machine) Return(o *Object) {
	f := VM.Frames[VM.CurFrame]
	f.Local[f.RetLocal] = o
}

func (VM *Machine) Lookup(name string) (o *Object) {
	if o = VM.Globals[name]; o != nil {
		return o
	}
	v := reflect.ValueOf(*VM.Builtins)
	fn := v.FieldByName("B_" + name)
	// fmt.Println("looked up", name, fn.Call([]reflect.Value{})[0].Interface().(*Object))
	if fn.IsValid() {
		return fn.Interface().(*Object)
	} else {
		panic("tried to lookup nonexistent name:  "+name)
	}
	return VM.Lookup("None")
}

func main() {
	 //builtin function call test case
	VM := NewMachine()
	consts := []*Object{VM.NewInt(3), VM.NewInt(2)}
	names := make([]string, 1)
	names[0] = "print"
	code := []Bytecode{
		{CONST, 0, 1, 0, 0},
		{CONST, 0, 2, 0, 1},
		{MUL, 0, 3, 1, 2},
		{ARGPUSH, 0, 3, LOCAL, 0},
		{GET, 0, 4, 0, 0},
		{CALL, 0, 0, 0, 4},
		{ARGPUSH, 0, 0, LOCAL, 0},
		{CALL, 0, 0, 0, 4}}
	 // int, raw_input, IF test case
	/*consts := []*Object{NewInt(3), NewString("less than 3"), NewString("equal to 3"), NewString("greater than 3")}
	names := make([]string, 3)
	names[0] = "raw_input"
	names[1] = "print"
	names[2] = "int"
	code := []Bytecode{
		{GET, 0, 1, 0, 0},
		{CALL, 0, 2, 1, 0},
		{GET, 0, 1, 0, 2},
		{ARGPUSH, 0, 2, LOCAL, 0},
		{CALL, 0, 2, 1, 0},
		{CONST, 0, 3, 0, 0},
		{LT, 0, 4, 2, 3},
		{IF, 0, 4, 0, 0},
		{JUMP, 0, 0, 4, 0},
		{ARGPUSH, 0, 1, CONSTANT, 0},
		{GET, 0, 5, 0, 1},
		{CALL, 0, 0, 5, 0},
		{JUMP, 0, 0, 10, 0},
		{EQ, 0, 4, 2, 3},
		{IF, 0, 4, 0, 0},
		{JUMP, 0, 0, 4, 0},
		{ARGPUSH, 0, 2, CONSTANT, 0},
		{GET, 0, 5, 0, 1},
		{CALL, 0, 0, 5, 0},
		{JUMP, 0, 0, 3, 0},
		{ARGPUSH, 0, 3, CONSTANT, 0},
		{GET, 0, 5, 0, 1},
		{CALL, 0, 0, 5, 0},
		{NOP, 0, 0, 0, 0}}*/
	// while loop test
	/*consts := []*Object{NewObject(10), NewObject(1)}
	names := make([]string, 1)
	names[0] = "print"
	code := []Bytecode{
		{CONST, 0, 1, 0, 0},
		{CONST, 0, 2, 0 , 1},
		{GET, 0, 3, 0, 0},
		{IF, 0, 1, 0, 0},
		{JUMP, 0, 0, 4, 0},
		{SUB, 0, 1, 1, 2},
		{ARGPUSH, 0, 1, LOCAL, 0},
		{CALL, 0, 0, 3, 0},
		{JUMP, 0, 0, 6, 1}}*/

	f := NewFrame(consts, names, code)
	VM.Run(f)
}
