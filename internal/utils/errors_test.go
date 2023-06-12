// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInternalError_Error(t *testing.T) {
	err := InternalError{Err: errors.New("Failed")}
	assert.Equal(t, err.Error(), "Internal Error: Failed")
}

func TestInternalError_Is(t *testing.T) {
	error1 := InternalError{}
	error2 := InternalError{Err: errors.New("Failed")}
	isInternalError := error1.Is(error2)
	assert.True(t, isInternalError)

	error3 := errors.New("Failed")
	isInternalError = error1.Is(error3)
	assert.False(t, isInternalError)
}

func TestSoftDeleteError_Is(t *testing.T) {
	error1 := SoftDeleteError(1)
	error2 := SoftDeleteError(1)
	isSoftDeleteError := error1.Is(error2)
	assert.True(t, isSoftDeleteError)
}

func TestNotFoundError_Is(t *testing.T) {
	error1 := NotFoundError(1)
	error2 := NotFoundError(1)
	isNotFoundError := error1.Is(error2)
	assert.True(t, isNotFoundError)
}

func TestInvalidAuthorizationError_Error(t *testing.T) {
	err := InvalidAuthorizationError{Err: errors.New("Auth Failed")}
	assert.Equal(t, err.Error(), "Invalid Authorization Error: Auth Failed")
}

func TestInvalidAuthorizationError_Is(t *testing.T) {
	error1 := InvalidAuthorizationError{}
	error2 := InvalidAuthorizationError{Err: errors.New("Auth Failed")}
	isInternalError := error1.Is(error2)
	assert.True(t, isInternalError)

	error3 := errors.New("Auth Failed")
	isInternalError = error1.Is(error3)
	assert.False(t, isInternalError)
}

func TestMissingAuthorizationError_Error(t *testing.T) {
	err := MissingAuthorizationError(-1)
	assert.Equal(t, err.Error(), MissingAuthorizationErrorMsg)
}

func TestNotAuthorizedError_Error(t *testing.T) {
	err := NotAuthorizedError(-1)
	assert.Equal(t, err.Error(), NotAuthorizedErrorMsg)
}

func TestRestClientError_Error(t *testing.T) {
	err := RestClientError{StatusCode: 201, Err: errors.New("Failed")}
	assert.Equal(t, err.Error(), "Failed")
}

func TestRestClientError_Is(t *testing.T) {
	error1 := RestClientError{}
	error2 := RestClientError{Err: errors.New("Failed")}
	isRestClientError := error1.Is(error2)
	assert.True(t, isRestClientError)

	error3 := errors.New("Failed")
	isRestClientError = error1.Is(error3)
	assert.False(t, isRestClientError)
}

func TestNotFoundError_Error(t *testing.T) {
	err := NotFoundError(-1)
	assert.Equal(t, err.Error(), NotFoundErrorMsg)
}

func TestSoftDeleteError_Error(t *testing.T) {
	err := SoftDeleteError(-1)
	assert.Equal(t, err.Error(), SoftDeleteErrorMsg)
}

func TestInvalidPathParamError_Error(t *testing.T) {
	err := InvalidPathParamError(-1)
	assert.Equal(t, err.Error(), InvalidPathParamErrorMsg)
}

func TestBodyNotFoundError_Error(t *testing.T) {
	err := BodyNotFoundError(-1)
	assert.Equal(t, err.Error(), BodyNotFoundErrorMsg)
}

func TestInvalidRequestFormatError_Error(t *testing.T) {
	err := InvalidRequestFormatError(-1)
	assert.Equal(t, err.Error(), InvalidRequestFormatErrorMsg)
}

func TestInvalidFilterError_Error(t *testing.T) {
	err := InvalidFilterError{Err: errors.New("Failed")}
	assert.Equal(t, err.Error(), "Invalid Filter Error: Failed")
}

func TestInvalidInputArgError_Error(t *testing.T) {
	err := InvalidInputArgError{Err: errors.New("Failed")}
	assert.Equal(t, err.Error(), "Invalid input argument Error: Failed")
}
