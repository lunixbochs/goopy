package main

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
)

func NewFunction(o NativeFunction) *Object {
	return NewObject(o);
}

type Builtins struct {}

func (g *Builtins) B_print() *Object {
	return NewFunction(func(VM *Machine, args chan *Object) {
		for len(args) > 0 {
			fmt.Print((<-args).Value, " ")
		}
		fmt.Println()
		cf := VM.Frames[VM.CurFrame]
		cf.Local[cf.RetLocal] = &Object{NONE, nil}
		return
	})
}

func (g *Builtins) B_raw_input() *Object {
	return NewFunction(func(VM *Machine, args chan *Object) {
		b := bufio.NewReader(os.Stdin)
		out, err := b.ReadBytes('\n')
		if err == nil {
			cf := VM.Frames[VM.CurFrame]
			cf.Local[cf.RetLocal] = NewObject(string(out[:len(out)-1]))
			return
		} else {
			panic("reading from stdin failed")
		}})
}

func (g *Builtins) B_int() *Object {
	return NewFunction(func(VM *Machine, args chan *Object) {
		s := (<-args).Value.(string)
		i, err := strconv.Atoi(s)
		if err == nil {
			cf := VM.Frames[VM.CurFrame]
			cf.Local[cf.RetLocal] = &Object{INT, i}
			return
		} else {
			fmt.Println(err)
			panic("need to raise an exception here instead")
		}})
}
