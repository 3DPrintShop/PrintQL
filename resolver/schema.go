package resolver

// The QueryResolver is the entry point for all top-level read operations.
type SchemaResolver struct {
}

func NewRoot() (*SchemaResolver, error) {
	return &SchemaResolver{}, nil
}