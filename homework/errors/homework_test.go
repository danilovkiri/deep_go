package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	var errText strings.Builder
	for _, subErr := range e.errs {
		errText.WriteString("\t* " + subErr.Error())
	}
	if errText.Len() > 0 {
		return strconv.Itoa(len(e.errs)) + " error(s) occurred:\n" + errText.String() + "\n"
	}

	return ""
}

func Append(err error, errs ...error) *MultiError {
	var newErr MultiError
	if err != nil {
		if mErr, ok := err.(*MultiError); ok {
			newErr.errs = append(newErr.errs, mErr.errs...)
		} else {
			newErr.errs = append(newErr.errs, err)
		}
	}
	newErr.errs = append(newErr.errs, errs...)
	return &newErr
}

func (e *MultiError) Unwrap() []error {
	return e.errs
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 error(s) occurred:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
