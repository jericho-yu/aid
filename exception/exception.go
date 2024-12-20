package exception

type (
	IException interface {
		Error() string
		Is(target error) bool
	}

	Exception struct {
		Err error
		Msg string
	}
)

func ExceptionError(e IException) {

}
