package domain

import "errors"

var (
	// data already exists error
	ErrDataAlreadyExists = errors.New("data already exists")

	// data not found error
	ErrDataNotFound = errors.New("data not found")

	// authentication error
	ErrAuthentication = errors.New("authentication failed")

	// authorization error
	ErrAuthToken = errors.New("invalid auth token")

	// internal server error
	ErrServer = errors.New("internal server error")

	// main category not found error
	ErrMainCategNotFound = errors.New("main category not found")

	// main category unique icon error
	ErrUniqueIconUser = errors.New("icon already used by another main category")

	// main category unique name error
	ErrUniqueNameUserType = errors.New("name already used by another main category with the same type")

	// sub category not found error
	ErrSubCategNotFound = errors.New("sub category not found")

	// sub category unique name error
	ErrUniqueNameUserMainCateg = errors.New("name already used by another sub category with the same main category")

	// icon not found error
	ErrIconNotFound = errors.New("icon not found")
)
