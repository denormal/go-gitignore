package gitignore

type Error interface {
	error
	Position() Position
} // Error()

type err struct {
	error
	_position Position
} // err()

func NewError(e error, p Position) Error {
	return &err{error: e, _position: p}
} // NewError()

func (e *err) Position() Position { return e._position }
