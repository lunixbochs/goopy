package main

func NewString(s string) *Object {
	math_proto := func(op func(a, b string) string) *Object{
		return NewFunction(func(VM *Machine, args []*Object) {
			a := (args[0]).Value.(string)
			b_o := args[1]
			var b string
			if b_o.Type == STRING {
				b = b_o.Value.(string)
			} else {
				panic("string op needs TypeError")
			}
			VM.Return(NewString(op(a, b)))
	})}
	attrs := make(map[string]*Object)
	attrs["__add__"] = math_proto(func(a, b string) string {return a+b})
	attrs["__mod__"] = math_proto(func(a, b string) string {
			panic("% based string formatting not implemented yet")
			return ""
		})
	attrs["__eq__"] = math_proto(func(a, b string) string {
			if a==b{ return a }
			return ""
		})
	return MakeObject(STRING, s, SubclassObject(), attrs)
}