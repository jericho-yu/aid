package myError

type (
	IMyError interface {
		New(msg string) IMyError
		Error() string
	}
	MyError struct{ msg string }
)
