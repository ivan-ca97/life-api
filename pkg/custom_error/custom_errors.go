package custom_error

import "log/slog"

type CustomError interface {
	error
	Severity() slog.Level
}

type HttpError interface {
	CustomError
	StatusCode() int
}

type CodedError interface {
	CustomError
	Code() string
}

type PublicMessageError interface {
	CustomError
	PublicMessage() string
}

type LoggableError interface {
	CustomError
	Log() map[string]string
}
