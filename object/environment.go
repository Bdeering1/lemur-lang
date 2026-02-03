package object

type Environment map[string]Object

func CreateEnvironment() *Environment {
    return &Environment{}
}
