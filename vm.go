package main

import (
	"reflect"
)

const (
	NAME = iota
	CONSTANT
	LOCAL
)

const (
	NONE = iota
	BOOL
	INT
	STRING
	LIST
	DICT
	FUNC
)

const (
	PRINT = iota
	RAW_INPUT
	INT_FUNC
)

const (
	NOP = iota
	ADD
	SUB 
	MUL 
	DIV 
	MOD 
	POW 
	BLSH
	BRSH
	BAND
	BOR 
	BXOR
	BNOT
	CMP
	NE 
	EQ 
	LE 
	LT 
	GET
	SET
	GETITEM
	SETITEM
	HASITEM
	DEL
	CONST
	IF
	JUMP
	ARGPUSH
	CALL
	CATCH
	RAISE
	RETURN
)


type Object struct {
	Type	int
	Value	interface{}
	//Attrs	map[string]*Object
}

type Func struct {
	name 	string
	Frame	*Frame
}
type NativeFunction func(VM *Machine, args chan *Object)

func NewObject(o interface{}) *Object {
	var typ int
	switch o.(type) {
		case bool: typ = BOOL
		case int: typ = INT
		case string: typ = STRING
		case []*Object: typ = LIST
		case map[*Object]*Object: typ = DICT
		case NativeFunction: typ = FUNC
		default: typ = NONE
	}

	return &Object{typ, o} //FIXME: init attrs
}

type Machine struct {
	Globals		map[string]*Object
	//Builtins	map[string]*Object
	Builtins 	*Builtins
	Frames		[]*Frame // keeps track of the call stack
	CurFrame 	int // starts at -1, ends at -1
}

func NewMachine() Machine {
	m := Machine{make(map[string]*Object), &Builtins{}, make([]*Frame, 0), -1}
	/*m.Builtins["print"] = &Object{FUNC, "print"}
	m.Builtins["raw_input"] = &Object{FUNC, RAW_INPUT}
	m.Builtins["int"] = &Object{FUNC, INT_FUNC}*/
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
	//VM.Builtins.print(VM, make(chan *Object))
	// fmt.Println(v.Kind(), v.MethodByName("Fprint"), VM, VM.Builtins)
	fn := v.MethodByName("B_" + name)
	if fn.IsValid() {
		return fn.Call([]reflect.Value{})[0].Interface().(*Object)
	} else {
		panic("tried to lookup nonexistent name:  "+name)
	}
	return &Object{NONE, nil}
}

type Bytecode struct {
	I 	byte
	A	byte
	B 	byte
	C 	byte
}

func (e *Bytecode) ABC() int {
	return int((e.A<<16)+(e.B<<8)+e.C)
}

func (e *Bytecode) BC() int {
	return int((e.B<<8)+e.C)
}

func (e *Bytecode) AB() int {
	return int((e.A<<8)+e.B)
}

type Frame struct {
	Const 		[]*Object
	Names 		[]string
	Local 		[]*Object
	Code 		[]Bytecode
	Cur 		int
	VM			*Machine
	Args 		chan *Object
	Ret 		*Object
	RetLocal	uint8
}

func NewFrame(Const []*Object, Names []string, Code []Bytecode) *Frame{
	f := &Frame{Const, Names, make([]*Object, 256), Code, 0, nil, make(chan *Object, 600), nil, 0}
	return f
}

func (f *Frame) Math(op Op) {
	e := f.Code[f.Cur]
	L := f.Local
	f.Local[e.A] = op(L[e.B], L[e.C])
	//fmt.Println("A, B and C are: ", L[e.A], L[e.B], L[e.C])
}

func (f *Frame) ArgPush(AB int, C byte) {
	// fmt.Println(AB, C, f.Local[AB])
	if len(f.Args) == cap(f.Args) {
		panic("not enough args for everyone")
	}
	switch (C) {
		case NAME: f.Args <- f.VM.Lookup(f.Names[AB])
		case CONSTANT: f.Args <- f.Const[AB]
		case LOCAL: f.Args <- f.Local[AB]
	}
}

func (f *Frame) Step() {
	e := f.Code[f.Cur]
	A, B, C := e.A, e.B, e.C
	L := f.Local
	VM := f.VM
	// fmt.Println(f.Cur, "e:  ", e, e.BC())
	switch (e.I) {
		case ADD: f.Math(Add)
		case SUB: f.Math(Sub)
		// case MUL: f.Math(Mul)
		// case DIV: f.Math(Div)
		// case MOD: f.Math(Mod)
		// case POW: f.Math(Pow)
		// case BLSH: f.Math(LeftShift)
		// case BRSH: f.Math(RightShift)
		// case BAND: f.Math(BitAnd)
		// case BOR: f.Math(BitOr)
		// case BXOR: f.Math(BitXor)
		// case BNOT: L[A] = BitNot(L[B])
		// case CMP: f.Math(Compare)
		// case NE: f.Math(NotEqual)
		case EQ: f.Math(Equal)
		//case LE: f.Math(LessEqual)
		case LT: f.Math(LessThan)
		case GET: L[A] = VM.Lookup(f.Names[e.BC()])
		/*case SET: VM.SetGlobal(f.Names[e.BC()], L[A])
		case GETITEM: VM.GetItem(L[A], L[B], L[C])
		case SETITEM: VM.SetItem(L[A], L[B], L[C])
		case HASITEM: VM.HasItem(L[A], L[B], L[C])*/
		case CONST: L[A] = f.Const[e.BC()]
		case IF: if Bool(L[A]) { f.Cur++ }
		case JUMP: if C != 0 { f.Cur -= e.AB() } else { f.Cur += e.AB() }
		case ARGPUSH: f.ArgPush(e.AB(), C)
		case CALL:
			f.RetLocal = A
			VM.Call(L[B], f.Args)
		case CATCH: break
		case RAISE: break
		case RETURN: f.Ret = L[A]
	}
	f.Cur++
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
	/*consts := []*Object{&Object{INT, 3}, &Object{STRING, "less than 3"}, &Object{STRING, "equal to 3"}, &Object{STRING, "greater than 3"}}
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
		{NOP, 0, 0, 0}}*/
	// while loop test
	consts := []*Object{&Object{INT, 10}, &Object{INT, 1}}
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
		{JUMP, 0, 6, 1}}
	VM := NewMachine()

	f := NewFrame(consts, names, code)
	VM.Run(f)
}
