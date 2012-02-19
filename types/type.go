package main

import (
	"fmt"
	//"reflect"
)

func (g *Builtins) B_object() *Object { // this represents python's 'object' type
	attrs := make(map[string]*Object)

	// FIXME: implement a None singleton and return it in lots of these calls
	attrs["__class__"] = NewString("type")
	attrs["__delattr__"] = NewFunction(func(VM *Machine, args []*Object){
			self := args[0]
			name := args[1].Value.(string)
			self.Attrs[name] = nil, false
		})
	attrs["__doc__"] = NewString("The most base type")
	attrs["__format__"] = NewFunction(func(VM *Machine, args []*Object){
			panic("object.__format__ not implemented!")
		})
	attrs["__getattribute__"] = NewFunction(func(VM *Machine, args []*Object){
			self := args[0]
			name := args[1].Value.(string)
			// fmt.Println("value", self.Attrs[name])
			VM.Return(self.Attrs[name])
		})
	attrs["__hash__"] = NewFunction(func(VM *Machine, args []*Object){
			self := args[0]
			// FIXME: !!!
			VM.Return(NewInt(self.Value.(int)))
			//VM.Return(NewInt(reflect.New(self).Addr().Interface().(int)))
		})
	attrs["__init__"] = NewFunction(func(VM *Machine, args []*Object){})
	attrs["__new__"] = NewFunction(func(VM *Machine, args []*Object){
			//FIXME: __new__ needs to actually do stuff
			typ := args[0]
			if typ.Type != TYPE {
				panic("can't create a new instance of a non-type object")
			}
		})
	attrs["__reduce__"] = NewFunction(func(VM *Machine, args []*Object){
			panic("pickling not implemented")
		})
	attrs["__reduce_ex__"] = NewFunction(func(VM *Machine, args []*Object){
			panic("pickling not implemented")
		})
	attrs["__repr__"] = NewFunction(func(VM *Machine, args []*Object){
			self := args[0]
			name := self.Value.(string)
			VM.Return(NewString(fmt.Sprintf("<type '%s'>\n", name)))
		})
	attrs["__setattr__"] = NewFunction(func(VM *Machine, args []*Object){
			self := args[0]
			name := args[1].Value.(string)
			value := args[2]
			self.Attrs[name] = value
		})
	attrs["__str__"] = NewFunction(func(VM *Machine, args []*Object){
			self := args[0]
			VM.Call(self.GetAttribute(VM, "__repr__"), args)
		})
	attrs["__subclasshook__"] = NewFunction(func(VM *Machine, args []*Object){
			self := args[0]
			bases := args[1].Bases
			for _, base := range bases {
				if base == self {
					VM.Return(NewBool(true))
				}
			}
			VM.Return(NewBool(false))
		})
	return MakeObject(TYPE, "object", nil, attrs)
}

func NewType(name string, bases []*Object, attrs map[string]*Object) *Object {
	return MakeObject(TYPE, name, bases, attrs)
}