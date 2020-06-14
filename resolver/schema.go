package resolver

// SchemaResolver is a root resolver for querying the graphql interface.
type SchemaResolver struct {
}

// NewRoot creates the root schema resolver.
func NewRoot() (*SchemaResolver, error) {
	return &SchemaResolver{}, nil
}
