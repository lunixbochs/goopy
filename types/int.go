package main

import (
	"fmt"
	"strconv"
)

func (g *Builtins) B_int() *Object {
	return NewFunction(func(VM *Machine, args []*Object) {
		s := (args[0]).Value.(string)
		i, err := strconv.Atoi(s)
		if err == nil {
			VM.Return(NewInt(i))
			return
		} else {
			fmt.Println(err)
			panic("need to raise an exception here instead")
		}})
}

func NewInt(i int) *Object {
	math_proto := func(op func(a, b int) int) *Object{
		return NewFunction(func(VM *Machine, args []*Object) {
			a := (args[0]).Value.(int)
			b_o := args[1]
			var b int
			switch b_o.Type { //FIXME: add support for floats etc.
				case INT: b = b_o.Value.(int)
				default: panic("need to raise a TypeError about invalid args")
			}
			res := NewInt(op(a, b))
			VM.Return(res)
	})}
	attrs := make(map[string]*Object)
	attrs["__add__"] = math_proto(func(a, b int) int {return a+b})
	attrs["__sub__"] = math_proto(func(a, b int) int {return a-b})
	attrs["__mul__"] = math_proto(func(a, b int) int {return a*b})
	attrs["__div__"] = math_proto(func(a, b int) int {return a/b})
	attrs["__mod__"] = math_proto(func(a, b int) int {return a%b})
	attrs["__lt__"] = math_proto(func(a, b int) int { if a<b {return 1}; return 0 })
	attrs["__le__"] = math_proto(func(a, b int) int { if a<=b {return 1}; return 0 })
	attrs["__eq__"] = math_proto(func(a, b int) int { if a==b {return 1}; return 0 })
	return MakeObject(INT, i, SubclassObject(), attrs)
}
