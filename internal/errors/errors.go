package ierrors

import "fmt"

type InfractionErrorKind int32

const (
	USER_ERR = iota
	SYS_ERR
)

type InfractionErrorI interface {
	error
	GetHttpCode() int
	GetCode() string
	GetMessage() string
	GetArgs() map[string]string
	GetDescription() string
	GetKind() InfractionErrorKind
	Unwrap() error
}
type InfractionError struct {
	Kind        InfractionErrorKind `json:"-"`
	Code        string              `json:"code"`
	HttpCode    int                 `json:"-"`
	Message     string              `json:"message"`
	Description string              `json:"description"`
	Cause       error               `json:"-"`
	Args        map[string]string   `json:"args"`
}

func NewInfractionErr(kind InfractionErrorKind, code string, httpCode int, message string, description string) *InfractionError {
	return &InfractionError{
		kind,
		code,
		httpCode,
		message,
		description,
		nil,
		nil,
	}
}

func (i *InfractionError) Wrap(err error) *InfractionError {
	i.Cause = err
	return i
}
func (i *InfractionError) Error() string {
	return fmt.Sprintf("%s", i.Message)
}
func (i *InfractionError) GetHttpCode() int {
	return i.HttpCode
}
func (i *InfractionError) GetCode() string {
	return i.Code
}
func (i *InfractionError) GetMessage() string {
	return i.Message
}
func (i *InfractionError) GetDescription() string {
	return i.Description
}

func (i *InfractionError) GetKind() InfractionErrorKind {
	return i.Kind
}
func (i *InfractionError) GetArgs() map[string]string {
	return i.Args
}

func (i *InfractionError) Unwrap() error {
	return i.Cause
}
