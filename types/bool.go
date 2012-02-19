package main

func Bool(a *Object) bool {
	switch a.Type {
		case BOOL: return a.Value.(bool)
		case INT: return a.Value.(int) != 0
		case NONE: return false
		case STRING: return len(a.Value.(string)) > 0
		case FUNC: return true
		case LIST: return len(a.Value.([]*Object)) > 0
		case DICT: return len(a.Value.(map[int]*Object)) > 0
		default: return true
	}
	panic("unable to convert to bool, contradiction achieved!")
}

func NewBool(a bool) *Object {
	// FIXME: support all the other stuff bools can do
	math_proto := func(op func(a, b bool) bool) *Object{
		return NewFunction(func(VM *Machine, args []*Object) {
			a := (args[0]).Value.(bool)
			b := Bool(args[1])
			VM.Return(NewBool(op(a, b)))
	})}
	attrs := make(map[string]*Object)
	attrs["__eq__"] = math_proto(func(a, b bool) bool {return a==b})
	return MakeObject(BOOL, a, SubclassObject(), attrs)
}