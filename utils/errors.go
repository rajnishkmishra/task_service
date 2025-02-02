package utils

import (
	"errors"
	"fmt"
	"net"
	"os"
	"runtime/debug"
	"syscall"

	"github.com/sirupsen/logrus"
)

type WrapperError interface {
	error
	BaseError() error
	HttpCode() int
	ErrCode() int
}

type simpleWrapperError struct {
	err     error
	errCode int
	status  int
}

func (s *simpleWrapperError) ErrCode() int {
	return s.errCode
}

func (s *simpleWrapperError) BaseError() error {
	return s.err
}

func (s *simpleWrapperError) HttpCode() int {
	return s.status
}

func (s *simpleWrapperError) Error() string {
	return s.err.Error()
}

func NewWrapperError(status int, err error) *simpleWrapperError {
	return &simpleWrapperError{
		status: status,
		err:    err,
	}
}

func Recovery() {
	if r := recover(); r != nil {
		if ne, ok := r.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if se.Err == syscall.EPIPE || se.Err == syscall.ECONNRESET {
					return
				}
			}
		}

		stackString := string(debug.Stack())
		fmt.Printf(stackString)
		err := errors.New(fmt.Sprintf("Recovered from following error:", r))
		logrus.WithError(err).Logln(logrus.ErrorLevel)
	}
}
