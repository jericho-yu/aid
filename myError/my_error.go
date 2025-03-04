package myError

type (
	IMyError interface {
		New(msg string) IMyError
		Warp(err error) IMyError
		Error() string
		Is(target error) bool
	}
	MyError struct{ Msg string }
)
