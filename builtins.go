package main

import (
	"fmt"
	"bufio"
	"os"
)

type Builtins struct {
	B_None		*Object
	B_NoneType	*Object
	B_int		*Object
	object		*Object
	B_print		*Object
}

var Builtin Builtins

func init() {
	Builtin.B_print = NewFunction(func(VM *Machine, args []*Object) {
		for _, x := range args {
			if x.Type != STRING {
				str := x.GetAttribute(VM, "__str__")
				fmt.Print(VM.RCall(str, []*Object{x}).Value.(string), " ")
			}
		}
		fmt.Println()
		cf := VM.Frames[VM.CurFrame]
		cf.Local[cf.RetLocal] = VM.Lookup("None")
		return
	})

}

// func (g *Builtins) B_print() *Object {
// 	fmt.Println("returning print")
// 	return NewFunction(func(VM *Machine, args []*Object) {
// 		for _, x := range args {
// 			str := x.GetAttribute(VM, "__str__")
// 			fmt.Println(str)
// 			if str != nil {
// 				fmt.Print(VM.RCall(str, []*Object{x}), " ")
// 			} else {
// 				fmt.Print(x.Value, " ")
// 			}
// 		}
// 		fmt.Println()
// 		cf := VM.Frames[VM.CurFrame]
// 		cf.Local[cf.RetLocal] = NewObject(nil)
// 		return
// 	})
// }

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

