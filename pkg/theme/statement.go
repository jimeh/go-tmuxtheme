package theme

type Statement interface {
	Parse(string) error
	Execute() error
}
