package main

import (
	// "fmt"
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
	I	byte
	A	byte
	B	byte
	C	byte
	D	byte
}
func (e *Bytecode) ABCD() int{
	return int(e.AB())<<16+int(e.CD())
}
func (e *Bytecode) AB() uint16{
	return (uint16(e.A)<<8)+uint16(e.B)
}

func (e *Bytecode) CD() uint16{
	return uint16((e.C<<8)+e.D)
}

type Frame struct {
	Const 		[]*Object
	Names 		[]string
	Local 		[]*Object
	Code 		[]Bytecode
	Cur 		int
	VM			*Machine
	Args 		[]*Object
	Ret 		*Object
	RetLocal	uint16
}

func FakeFrame(VM *Machine) *Frame{
	f := &Frame{make([]*Object, 0), make([]string, 0), make([]*Object, 65536), make([]Bytecode, 0), 0, VM, make([]*Object, 0), nil, 0}
	return f
}

func NewFrame(Const []*Object, Names []string, Code []Bytecode) *Frame{
	f := &Frame{Const, Names, make([]*Object, 65536), Code, 0, nil, make([]*Object, 0), nil, 0}
	return f
}

func (f *Frame) Math(op string) {
	e := f.Code[f.Cur]
	//push the operation's arguments and set the return location
	f.RetLocal = e.AB()
	f.Args = append(f.Args, f.Local[e.C])
	f.Args = append(f.Args, f.Local[e.D])
	//lookup the actual operation
	//FIXME: on failure, try the __r*__ variant? raise exception?
	opfunc := f.Local[e.C].GetAttribute(f.VM, op)
	// fmt.Println("preparing to call", op)
	f.VM.Call(opfunc, f.Args)
	//fmt.Println("A, B and C are: ", L[e.A], L[e.B], L[e.C])
}

func (f *Frame) ArgPush(AB uint16, C byte) {
	// fmt.Println(AB, C, f.Local[AB])
	switch (C) {
		case NAME: f.Args = append(f.Args, f.VM.Lookup(f.Names[AB]))
		case CONSTANT: f.Args = append(f.Args, f.Const[AB])
		case LOCAL: f.Args = append(f.Args, f.Local[AB])
	}
}

func (f *Frame) Step() {
	e := f.Code[f.Cur]
	A, C, AB, CD, ABCD := e.A, e.C, e.AB(), e.CD(), e.ABCD()
	L := f.Local
	VM := f.VM
	// fmt.Println(f.Cur, "e:  ", e, AB, CD, f.Local[:6])
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
		case GET: L[AB] = VM.Lookup(f.Names[CD])
		/*case SET: VM.SetGlobal(f.Names[e.BC()], L[A])
		case GETITEM: VM.GetItem(L[A], L[B], L[C])
		case SETITEM: VM.SetItem(L[A], L[B], L[C])
		case HASITEM: VM.HasItem(L[A], L[B], L[C])*/
		case CONST: L[AB] = f.Const[CD]
		case IF: if Bool(L[A]) { f.Cur++ }
		case JUMP: f.Cur += ABCD // TODO: make signed
		case ARGPUSH: f.ArgPush(AB, C)
		case CALL:
			f.RetLocal = AB
			VM.Call(L[CD], f.Args)
			f.Args = make([]*Object, 0)
		case CATCH: break
		case RAISE: break
		case RETURN: f.Ret = L[AB]
	}
	f.Cur++
}