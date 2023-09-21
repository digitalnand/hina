package hina

type Environment struct {
	SymbolTable map[string]interface{}
}

func NewEnvironment() Environment {
	var env Environment
	env.SymbolTable = make(map[string]interface{})
	return env
}

func (env Environment) Get(identifier string) (any, bool) {
	node, hasNode := env.SymbolTable[identifier]
	if !hasNode {
		return nil, hasNode
	}
	return node, hasNode
}

func (env Environment) Set(identifier string, value any) {
	env.SymbolTable[identifier] = value
}
