package cmpopts

// Options is a stub type to satisfy imports from dependencies.
type Options struct{}

// IgnoreFields returns an empty Options value.
func IgnoreFields(_ interface{}, _ ...string) Options {
	return Options{}
}
