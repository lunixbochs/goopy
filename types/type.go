package main

import (
	"fmt"
	//"reflect"
)

func (g *Builtins) B_object() *Object { // this represents python's 'object' type
	return SubclassObject()[0]
}

func NewObjectRoot(ro *Object) *Object{
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
			for _, base := range append([]*Object{self}, self.Bases...) {
				if base.Attrs[name] != nil {
					VM.Return(base.Attrs[name])
					return
				}
			}
			VM.Return(self.Attrs[name])
		})
	attrs["__hash__"] = NewFunction(func(VM *Machine, args []*Object){
			self := args[0]
			// FIXME: !!!
			VM.Return(VM.NewInt(self.Value.(int)))
			//VM.Return(NewInt(reflect.New(self).Addr().Interface().(int)))
		})
	attrs["__init__"] = NewFunction(func(VM *Machine, args []*Object){})
	attrs["__new__"] = NewFunction(func(VM *Machine, args []*Object){
			class := args[0]
			iargs := args[1:]
			obj := MakeObject(OBJECT, 234, class.Bases, class.Attrs)
			iargs = append([]*Object{obj}, iargs...)
			// __init__ must not have a return value
			VM.RCall(obj.GetAttribute(VM, "__init__"), iargs)
			VM.Return(obj)
		})
	attrs["__reduce__"] = NewFunction(func(VM *Machine, args []*Object){
			panic("pickling not implemented")
		})
	attrs["__reduce_ex__"] = NewFunction(func(VM *Machine, args []*Object){
			panic("pickling not implemented")
		})
	attrs["__repr__"] = NewFunction(func(VM *Machine, args []*Object){
			self := args[0]
			name := self.Value
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
	ro.Type = TYPE
	ro.Value = "object"
	ro.Attrs = attrs
	return ro
}

func NewType(name string, bases []*Object, attrs map[string]*Object) *Object {
	return MakeObject(TYPE, name, bases, attrs)
}