package cfgds

// CtorCtor is a function that creates new Ctors
type CtorCtor func() Ctor

var registeredCtors = map[string]CtorCtor{
	"mem": func() Ctor { return &memCtor{} },
}

// RegisterCtor allows for registration of Ctor under given name/type
func RegisterCtor(name string, cctor CtorCtor) {
	registeredCtors[name] = cctor
}
