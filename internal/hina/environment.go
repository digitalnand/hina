package hina

type Environment struct {
	SymbolTable map[string]interface{}
}

func NewEnvironment() Environment {
	var env Environment
	env.SymbolTable = make(map[string]interface{})
	return env
}

func (env Environment) Get(identifier string) (Term, bool) {
	node, hasNode := env.SymbolTable[identifier]
	if !hasNode {
		return nil, false
	}
	return node, true
}

func (env Environment) Set(identifier string, value Term) {
	env.SymbolTable[identifier] = value
}
