package utils

// Must evaluate the given expression and panic if it fails.
func Must[T any](expr T, err error) T {
	if err != nil {
		panic(err)
	}
	return expr
}

// MustWithMessage evaluate the given expression and panic if it fails.
func MustWithMessage[T any](expr T, err error, msg string) T {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
	return expr
}
