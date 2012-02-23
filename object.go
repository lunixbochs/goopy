package main

import (
	// "fmt"
)

const ( //enumerate the builtin types
	NONE = iota
	BOOL
	INT
	STRING
	LIST
	DICT
	FUNC
	TYPE
)

type Object struct {
	Type	int
	Value	interface{}
	Attrs	map[string]*Object
	Bases	[]*Object
}

type Func struct {
	name 	string
	Frame	*Frame
}

type NativeFunction func(VM *Machine, args []*Object)

func NewFunction(o NativeFunction) *Object {
	return NewObject(o);
}

func MakeObject(typ int, o interface{}, bases []*Object, attrs map[string]*Object)*Object {
	return &Object{typ, o, attrs, bases}
}

func NewObject(o interface{}) *Object {
	var typ int
	switch o.(type) {
		case bool: typ = BOOL
		case int: return NewInt(o.(int))
		case string: typ = STRING
		case []*Object: typ = LIST
		case map[*Object]*Object: typ = DICT
		case NativeFunction: typ = FUNC
		case nil: typ = NONE
		default: typ = NONE
	}

	return MakeObject(typ, o, make([]*Object, 0), make(map[string]*Object))	
}

func (o *Object) GetAttribute(VM *Machine, name string) *Object {
	for _, base := range o.Bases {
		a := make([]*Object, 0)
		a = append(a, o)
		a = append(a, NewString(name))
		// fmt.Println("about to RCall", *a[0])
		ret := VM.RCall(base.Attrs["__getattribute__"], a)
		if ret != nil {
			return ret
		}
	}
	panic("attribute not found!")
}

var real_object *Object = nil

func SubclassObject() []*Object{
	if real_object == nil {
		real_object = &Object{}
		real_object = NewObjectRoot(real_object)
		// fmt.Println(real_object)
	}
	return []*Object{real_object}
}