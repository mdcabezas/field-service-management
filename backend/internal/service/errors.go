package service

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation: %s: %s", e.Field, e.Message)
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s %s", e.Resource, e.ID)
}

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s", e.Message)
}

type TransitionError struct {
	Entity string
	From   string
	To     string
}

func (e *TransitionError) Error() string {
	return fmt.Sprintf("invalid transition: %s from %s to %s", e.Entity, e.From, e.To)
}

// HandleRepoGetByIDError wraps errors from repo GetByID calls.
// If the error is already a NotFoundError, it returns a new NotFoundError with the given resource/id.
// Otherwise, it wraps the original error with context so DB connection failures surface as 500.
func HandleRepoGetByIDError(err error, resource, id string) error {
	if err == nil {
		return nil
	}
	var nf *NotFoundError
	if errors.As(err, &nf) {
		return &NotFoundError{Resource: resource, ID: id}
	}
	return fmt.Errorf("get %s by id %s: %w", resource, id, err)
}
