package util

import "fmt"

type ExistingRecordError struct {
	Identification string
}

func (e *ExistingRecordError) Error() string {
	return fmt.Sprintf("A record by [%s] already exists in the database!", e.Identification)
}

type NotExistingRecordError struct {
	Identification string
}

func (e *NotExistingRecordError) Error() string {
	return fmt.Sprintf("A record by [%s] does not exist in the database!", e.Identification)
}
