package main

import (
	"fmt"
	"bufio"
	"os"
)

type Builtins struct {
	none	*Object
	object	*Object
}

func (g *Builtins) B_print() *Object {
	return NewFunction(func(VM *Machine, args []*Object) {
		for _, x := range args {
			fmt.Print(x.Value, " ")
		}
		fmt.Println()
		cf := VM.Frames[VM.CurFrame]
		cf.Local[cf.RetLocal] = NewObject(nil)
		return
	})
}

func (g *Builtins) b_None() *Object {
	if g.none == nil {
		g.none = NewObject(nil)
	}
	return g.none
}

func (g *Builtins) B_raw_input() *Object {
	return NewFunction(func(VM *Machine, args []*Object) {
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

