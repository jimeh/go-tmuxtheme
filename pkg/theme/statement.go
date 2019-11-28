package theme

type Statement interface {
	Parse(string) error
	Execute(theme *Theme) error
}

func NewStatement(body string) (Statement, error) {
	statements := []Statement{
		&EmptyStatement{},
		&CommentStatement{},
		&SetOptionStatement{},
	}

	for _, t := range statements {
		err := t.Parse(body)
		if err == nil {
			return t, nil
		}
		if _, ok := err.(*NotSupportedCommandError); !ok {
			return nil, err
		}
	}

	return nil, &UnsupportedStatementError{Body: body}
}
