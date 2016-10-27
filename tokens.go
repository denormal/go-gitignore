package gitignore

type Tokens []*Token

func (t Tokens) String() string {
	// concatenate the tokens into a single string
	_rtn := ""
	for _, _t := range []*Token(t) {
		_rtn = _rtn + _t.Token()
	}
	return _rtn
} // String()
