package theme

type Theme struct {
	ServerOptions        map[string]string
	GlobalSessionOptions map[string]string
	SessionOptions       map[string]string
	GlobalWindowOptions  map[string]string
	WindowOptions        map[string]string
	Statements           []*Statement
}

func NewTheme() *Theme {
	return &Theme{
		ServerOptions:        map[string]string{},
		GlobalSessionOptions: map[string]string{},
		SessionOptions:       map[string]string{},
		GlobalWindowOptions:  map[string]string{},
		WindowOptions:        map[string]string{},
		Statements:           []*Statement{},
	}
}
