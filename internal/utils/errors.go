// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package utils

import "fmt"

// Error Messages
const (
	MissingAuthorizationErrorMsg = "User Authorization information is missing"
	NotAuthorizedErrorMsg        = "User is not authorized"
	InvalidAuthorizationErrorMsg = "Invalid Authorization Error"
	NotFoundErrorMsg             = "Object Not Found"
	SoftDeleteErrorMsg           = "Object is Deleted"
	InternalServerErrorMsg       = "An internal server error occurred"
	BodyNotFoundErrorMsg         = "Request body not found"
	InvalidPathParamErrorMsg     = "Path parameter is invalid"
	InvalidRequestFormatErrorMsg = "Request format is invalid"
)

// InternalError: All kinds of connection errors (to DB or kafka), marshaling errors
// are mapped to Internal Errors
type InternalError struct {
	Err error
}

func (e InternalError) Error() string {
	return fmt.Sprintf("Internal Error: %s", e.Err.Error())
}

func (e InternalError) Is(target error) bool {
	if _, ok := target.(InternalError); ok {
		return true
	}
	return false
}

// Authorization Errors

type InvalidAuthorizationError struct {
	Err error
}

func (e InvalidAuthorizationError) Error() string {
	return fmt.Sprintf("%s: %s", InvalidAuthorizationErrorMsg, e.Err.Error())
}

func (e InvalidAuthorizationError) Is(target error) bool {
	if _, ok := target.(InvalidAuthorizationError); ok {
		return true
	}
	return false
}

type MissingAuthorizationError int

func (e MissingAuthorizationError) Error() string {
	return MissingAuthorizationErrorMsg
}

type NotAuthorizedError int

func (e NotAuthorizedError) Error() string {
	return NotAuthorizedErrorMsg
}

// Special Errors -
// NotFoundError - If GET by ID couldn't find any record
// SoftDeleteError - If GET by ID found a deleted record
// InvalidFilterError - 1. Filter from GET All request is not supported
//                      2. REST call URL has unsupported filters/operators

type NotFoundError int

func (e NotFoundError) Error() string {
	return NotFoundErrorMsg
}

func (e NotFoundError) Is(target error) bool {
	if _, ok := target.(NotFoundError); ok {
		return true
	}
	return false
}

type SoftDeleteError int

func (e SoftDeleteError) Error() string {
	return SoftDeleteErrorMsg
}

func (e SoftDeleteError) Is(target error) bool {
	if _, ok := target.(SoftDeleteError); ok {
		return true
	}
	return false
}

type InvalidFilterError struct {
	Err error
}

func (e InvalidFilterError) Error() string {
	return fmt.Sprintf("Invalid Filter Error: %s", e.Err.Error())
}

// Invalid input argument error

type InvalidInputArgError struct {
	Err error
}

func (e InvalidInputArgError) Error() string {
	return fmt.Sprintf("Invalid input argument Error: %s", e.Err.Error())
}

// Invalid path parameters error

type InvalidPathParamError int

func (e InvalidPathParamError) Error() string {
	return InvalidPathParamErrorMsg
}

// Request body not found error

type BodyNotFoundError int

func (e BodyNotFoundError) Error() string {
	return BodyNotFoundErrorMsg
}

// Invalid request format error

type InvalidRequestFormatError int

func (e InvalidRequestFormatError) Error() string {
	return InvalidRequestFormatErrorMsg
}

// REST client error

type RestClientError struct {
	StatusCode int
	Err        error
}

func (e RestClientError) Error() string {
	return e.Err.Error()
}

func (e RestClientError) Is(target error) bool {
	if _, ok := target.(RestClientError); ok {
		return true
	}
	return false
}
