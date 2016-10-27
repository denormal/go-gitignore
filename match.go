package gitignore

// Match represents the interface of successful matches against a .gitignore
// pattern set. A Match can be queried to determine whether the matched path
// should be ignored or included (i.e. was the path matched by a negated
// pattern), and to extract the position of the pattern within the .gitignore,
// and a string representation of the pattern.
type Match interface {
	Ignore() bool
	Include() bool
	String() string
	Position() Position
} // Match{}
