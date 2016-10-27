package gitignore

type Match interface {
	Ignore() bool
	Accept() bool
	String() string
	Position() Position
} // Match{}
