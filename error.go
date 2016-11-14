package gitignore

// Error is the interface for errors
type Error interface {
	error
	Position() Position
	Is(error) bool
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

// Is returns true if the Error instance is the same as the given error. This
// permits direct comparison against package errors such as CarriageReturnError.
func (e *err) Is(er error) bool {
	return e.error == er
} // Is()
