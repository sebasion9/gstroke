package errors

import "fmt"

type BaseDecoderError struct {
	Message string
}

func(e* BaseDecoderError) Error() string {
	return fmt.Sprintf("[ERR] %s\n", e.Message)
}

type InvalidJPEGError struct {
	*BaseDecoderError
}

func NewInvalidJPEGError(msg string) InvalidJPEGError {
	return InvalidJPEGError{BaseDecoderError: &BaseDecoderError{msg}}
}
