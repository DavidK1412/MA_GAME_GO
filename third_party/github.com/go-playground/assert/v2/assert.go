package assert

// Equal is a minimal stub used to satisfy dependency resolution without pulling the original module.
func Equal(expected, actual interface{}) bool {
	return expected == actual
}
