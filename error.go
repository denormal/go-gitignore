package gitignore

// Error is the interface for errors
type Error interface {
	error
	Position() Position
} // Error()

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
