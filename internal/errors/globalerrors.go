package ierrors

var (
	ErrNoInfractionFound = NewInfractionErr(USER_ERR, "1000_NOT_FOUND", 404, "the infraction could not be found", "")
	ErrValidationFailed  = NewInfractionErr(USER_ERR, "1001_VALIDATON_FAILED", 400, "validation failed", "")
	ErrDbFailure         = NewInfractionErr(SYS_ERR, "2000_INTERNAL", 500, "something went wrong", "please try again later")
	ErrUnspecified       = NewInfractionErr(USER_ERR, "2001_INTERNAL", 500, "something went wrong", "please try again later")
)
