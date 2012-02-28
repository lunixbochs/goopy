package main

import (
	"fmt"
	"strconv"
)

func init() {
	math_proto := func(op func(a, b int) int) *Object{
		return NewFunction(func(VM *Machine, args []*Object) {
			a := (args[0]).Value.(int)
			b_o := args[1]
			var b int
			switch b_o.Type { //FIXME: add support for floats etc.
				case INT: b = b_o.Value.(int)
				default: panic("need to raise a TypeError about invalid args")
			}
			res := VM.NewInt(op(a, b))
			VM.Return(res)
	})}
	attrs := make(map[string]*Object)
	attrs["__init__"] = NewFunction(func(VM *Machine, args []*Object){
			s := args[0]
			value := args[1].Value
			switch value.(type) {
				case int:
					s.Value = value.(int)
				case string:
					i, err := strconv.Atoi(value.(string))
					if err == nil {
						s.Value = i
					} else {
						fmt.Println(err)
						panic("need to raise an exception here instead")
					}
				// TODO: add more cases
			}
			s.Type = INT
		})
	attrs["__add__"] = math_proto(func(a, b int) int {return a+b})
	attrs["__sub__"] = math_proto(func(a, b int) int {return a-b})
	attrs["__mul__"] = math_proto(func(a, b int) int {return a*b})
	attrs["__div__"] = math_proto(func(a, b int) int {return a/b})
	attrs["__mod__"] = math_proto(func(a, b int) int {return a%b})
	attrs["__lt__"] = math_proto(func(a, b int) int { if a<b {return 1}; return 0 })
	attrs["__le__"] = math_proto(func(a, b int) int { if a<=b {return 1}; return 0 })
	attrs["__eq__"] = math_proto(func(a, b int) int { if a==b {return 1}; return 0 })
	attrs["__cmp__"] = math_proto(func(a, b int) int { if a>b {return 1} else if a==b {return 0}; return -1 })
	attrs["__str__"] = NewFunction(func(VM *Machine, args []*Object){
			VM.Return(NewString(strconv.Itoa(args[0].Value.(int))))
		})
	Builtin.B_int = MakeObject(TYPE, "int", SubclassObject(), attrs)
}

func (VM *Machine) NewInt(i int) *Object{
	return VM.RCall(VM.Lookup("int").GetAttribute(VM, "__new__"), []*Object{VM.Lookup("int"), Container(i)})
}
