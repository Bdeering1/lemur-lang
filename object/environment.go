package object

type Environment struct {
    store map[string]Object
    outer *Environment
}

func CreateEnvironment() *Environment {
    return &Environment{
        store: make(map[string]Object),
        outer: nil,
    }
}

func CreateEnclosedEnvironment(outer *Environment) *Environment {
    return &Environment{
        store: make(map[string]Object),
        outer: outer,
    }
}

func (e *Environment) Get(key string) (Object, bool) {
    obj, ok := e.store[key]
    if !ok && e.outer != nil { return e.outer.Get(key) }

    return obj, ok
}

func (e *Environment) Set(key string, val Object) { e.store[key] = val }
