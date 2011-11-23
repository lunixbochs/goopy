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
	Attrs	map[string]*Object
}

type Func struct {
	name 	string
	Frame	*Frame
}

type NativeFunction func(VM *Machine, args chan *Object)

func NewFunction(o NativeFunction) *Object {
	return NewObject(o);
}

func MakeObject(typ int, o interface{}, attrs map[string]*Object)*Object {
	return &Object{typ, o, attrs}
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

	return MakeObject(typ, o, make(map[string]*Object))	
}

func Bool(a *Object) bool {
	switch a.Type {
		case BOOL: return a.Value.(bool)
		case INT: return a.Value.(int) != 0
		case NONE: return false
		case STRING: return len(a.Value.(string)) > 0
		case FUNC: return true
		case LIST: return len(a.Value.([]*Object)) > 0
		case DICT: return len(a.Value.(map[int]*Object)) > 0
	}
	return false
}

func NewInt(i int) *Object {
	math_proto := func(op func(a, b int) int) *Object{
		return NewFunction(func(VM *Machine, args chan *Object) {
			a := (<-args).Value.(int)
			b_o := <-args
			var b int
			switch b_o.Type { //FIXME: add support for floats etc.
				case INT: b = b_o.Value.(int)
				default: panic("need to raise a TypeError about invalid args")
			}
			cf := VM.Frames[VM.CurFrame]
			cf.Local[cf.RetLocal] = NewInt(op(a, b))
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
	return MakeObject(INT, i, attrs)
}