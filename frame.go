package main

import (
	"fmt"
)

const (
	NAME = iota
	CONSTANT
	LOCAL
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

func (f *Frame) Math(op string) {
	e := f.Code[f.Cur]
	//push the operation's arguments and set the return location
	f.RetLocal = e.A
	f.Args <- f.Local[e.B]
	f.Args <- f.Local[e.C]
	//lookup the actual operation
	//FIXME: on failure, try the __r*__ variant? raise exception?
	opfunc, exists := f.Local[e.B].Attrs[op]
	fmt.Println("preparing to call", op)
	if !exists {
		panic("no such operation, needs to raise TypeError")
	}
	f.VM.Call(opfunc, f.Args)
	//fmt.Println("A, B and C are: ", L[e.A], L[e.B], L[e.C])
}

func (f *Frame) ArgPush(AB int, C byte) {
	// fmt.Println(AB, C, f.Local[AB])
	if len(f.Args) == cap(f.Args) {
		//FIXME: dynamically make new channels with twice the size, or maybe use vectors?
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
		case ADD: f.Math("__add__")
		case SUB: f.Math("__sub__")
		case MUL: f.Math("__mul__")
		case DIV: f.Math("__div__")
		case MOD: f.Math("__mod__")
		case POW: f.Math("__pow__")
		// case BLSH: f.Math(LeftShift)
		// case BRSH: f.Math(RightShift)
		// case BAND: f.Math(BitAnd)
		// case BOR: f.Math(BitOr)
		// case BXOR: f.Math(BitXor)
		// case BNOT: L[A] = BitNot(L[B])
		// case CMP: f.Math(Compare)
		// case NE: f.Math(NotEqual)
		case EQ: f.Math("__eq__")
		case LE: f.Math("__le__")
		case LT: f.Math("__lt__")
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