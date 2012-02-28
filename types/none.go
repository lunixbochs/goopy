package main

import (
	// "fmt"
)

func init() {
	attrs := make(map[string]*Object)
	attrs["__init__"] = NewFunction(func(VM *Machine, args []*Object){
			s := args[0]
			s.Type = NONE
			s.Value = nil
		})
	attrs["__str__"] = NewFunction(func(VM *Machine, args []*Object){
			VM.Return(NewString("None"))
		})
	Builtin.B_NoneType = MakeObject(TYPE, "NoneType", SubclassObject(), attrs)
}

func (VM *Machine) NewNone() *Object {
	return VM.RCall(VM.Lookup("NoneType").GetAttribute(VM, "__new__"), []*Object{VM.Lookup("NoneType")})
}
