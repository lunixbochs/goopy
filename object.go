package main

const (
	NONE = iota
	BOOL
	INT
	STRING
	LIST
	DICT
	FUNC
)

type Object struct {
	Type	int
	Value	interface{}
	//Attrs	map[string]*Object
}

type Func struct {
	name 	string
	Frame	*Frame
}

type NativeFunction func(VM *Machine, args chan *Object)

func NewFunction(o NativeFunction) *Object {
	return NewObject(o);
}

func NewObject(o interface{}) *Object {
	var typ int
	switch o.(type) {
		case bool: typ = BOOL
		case int: typ = INT
		case string: typ = STRING
		case []*Object: typ = LIST
		case map[*Object]*Object: typ = DICT
		case NativeFunction: typ = FUNC
		default: typ = NONE
	}

	return &Object{typ, o} //FIXME: init attrs
}
