package errors

import "errors"

const (
	ErrorDB           = "Database Error"
	ErrorExistingData = "Data Already Exist"
	ErrorNoDataChange = "No Data Changes"
	ErrorDataNotFound = "Data Not Found"
	ErrorUserNotFound = "User Not Found"
	ErrorInvalidUser  = "Invalid User"
)

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("Your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("Your Item already exist")
	// ErrInvalidRequest will throw if the given request-body or params is not valid
	ErrInvalidRequest = errors.New("Invalid Request")
	// ErrFailAuth will throw if the auth is not valid
	ErrFailAuth = errors.New("Authentication Failed")
)
