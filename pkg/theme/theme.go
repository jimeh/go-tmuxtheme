package theme

import (
	"bufio"
	"io"
	"os"
)

type Theme struct {
	ServerOptions        map[string]string
	GlobalSessionOptions map[string]string
	SessionOptions       map[string]string
	GlobalWindowOptions  map[string]string
	WindowOptions        map[string]string
	Statements           []Statement
}

func NewTheme() *Theme {
	return &Theme{
		ServerOptions:        map[string]string{},
		GlobalSessionOptions: map[string]string{},
		SessionOptions:       map[string]string{},
		GlobalWindowOptions:  map[string]string{},
		WindowOptions:        map[string]string{},
		Statements:           []Statement{},
	}
}

func (s *Theme) Parse(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	line := []byte{}

	for scanner.Scan() {
		line = append(line, scanner.Bytes()...)
		if len(line) > 0 && line[len(line)-1] == byte('\\') {
			line = line[:len(line)-1]
		} else {
			statement, err := NewStatement(string(line))
			if err != nil {
				return err
			}

			s.Statements = append(s.Statements, statement)
			line = []byte{}
		}
	}

	return scanner.Err()
}

func (s *Theme) Execute() error {
	for _, st := range s.Statements {
		err := st.Execute(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Theme) Load(filename string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	return s.Parse(r)
}
