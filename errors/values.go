package errors

var (
	// Unable to resolve is the standard message given when a value can't be resolved within the graphql schema.
	UnableToResolve = New("unable to resolve")
)

// WrongType creates anb error with formating for when the type of an interface isn't what was expected.
func WrongType(expected, actual interface{}) error {
	return Errorf("wrong type: wanted %T, got %T", expected, actual)
}
