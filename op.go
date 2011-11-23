package main

type Op func(*Object, *Object) *Object

func Add(a, b *Object) *Object {
	return &Object{INT, a.Value.(int) + b.Value.(int)}
}

func Sub(a, b *Object) *Object {
	return &Object{INT, a.Value.(int) - b.Value.(int)}
}

func Equal(a, b *Object) *Object {
	return &Object{BOOL, a.Value.(int) == b.Value.(int)}
}

func LessThan(a, b *Object) *Object {
	return &Object{BOOL, a.Value.(int) < b.Value.(int)}
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
