package gitignore

type Error interface {
	error
	Position() Position
	Underlying() error
}

// err extends the standard error to include a Position within the parsed
// .gitignore file
type err struct {
	error
	_position Position
} // err()

// NewError returns a new Error instance for the given error e and position p.
func NewError(e error, p Position) Error {
	return &err{error: e, _position: p}
} // NewError()

// Position returns the position of the error (i.e. the location within the
// .gitignore file)
func (e *err) Position() Position { return e._position }

// Underlying returns the underlying error, permitting direct comparison
// against the wrapped error.
func (e *err) Underlying() error {
	return e.error
} // Underlying()

// ensure err satisfies the Error interface
var _ Error = &err{}
