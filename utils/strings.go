package utils

func FistNonEmptyString(args ...*string) *string {
	for _, arg := range args {
		if arg != nil && len(*arg) > 0 {
			return arg
		}
	}
	return nil
}

// StringSlice returns the first non-nil value in a list of slice of string.
func StringSlice(args ...[]string) []string {
	for _, arg := range args {
		if arg != nil {
			return arg
		}
	}
	return nil
}
