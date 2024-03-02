package efactory

import "errors"

var (
	errInsertAssWithoutAss       = errors.New("inserting associations without any associations")
	errDestValueNotStruct        = errors.New("destination value is not a struct")
	errSourceValueNotStruct      = errors.New("source value is not a struct")
	errDestAndSourceIsDiff       = errors.New("destination and source type is different")
	errBuildListNGreaterThanZero = errors.New("BuildList: n must be greater than 0")
	errDBNotProvided             = errors.New("DB is not provided")
	errTableNameNotProvided      = errors.New("TableName is not provided")
)
