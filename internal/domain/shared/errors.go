// Package shared contains domain primitives used by every module: identifiers,
// money, civil dates, clocks and typed domain errors. It has no dependencies
// outside the standard library.
package shared

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Code is a stable, machine-readable error code shared by all transports.
type Code string

const (
	CodeValidation      Code = "validation"
	CodeNotFound        Code = "not_found"
	CodeConflict        Code = "conflict"
	CodeForbidden       Code = "forbidden"
	CodeUnauthenticated Code = "unauthenticated"
	CodeInvariant       Code = "invariant_violation"
	CodeLocked          Code = "account_locked"
	CodeInternal        Code = "internal"
)

// Error is the typed domain/application error. Adapters map Code to transport
// specific representations (HTTP status, Wails error envelope).
type Error struct {
	Code    Code
	Message string
	Fields  map[string]string
	Err     error
}

func (e *Error) Error() string {
	var b strings.Builder
	b.WriteString(string(e.Code))
	if e.Message != "" {
		b.WriteString(": ")
		b.WriteString(e.Message)
	}
	if len(e.Fields) > 0 {
		keys := make([]string, 0, len(e.Fields))
		for k := range e.Fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		b.WriteString(" [")
		for i, k := range keys {
			if i > 0 {
				b.WriteString("; ")
			}
			fmt.Fprintf(&b, "%s: %s", k, e.Fields[k])
		}
		b.WriteString("]")
	}
	if e.Err != nil {
		b.WriteString(": ")
		b.WriteString(e.Err.Error())
	}
	return b.String()
}

func (e *Error) Unwrap() error { return e.Err }

// Is allows errors.Is(err, shared.ErrNotFound) style comparisons by code.
func (e *Error) Is(target error) bool {
	var t *Error
	if errors.As(target, &t) {
		return t.Code == e.Code && t.Message == "" && t.Fields == nil
	}
	return false
}

// Sentinel errors for errors.Is comparisons.
var (
	ErrNotFound        = &Error{Code: CodeNotFound}
	ErrConflict        = &Error{Code: CodeConflict}
	ErrForbidden       = &Error{Code: CodeForbidden}
	ErrUnauthenticated = &Error{Code: CodeUnauthenticated}
	ErrInvariant       = &Error{Code: CodeInvariant}
	ErrValidation      = &Error{Code: CodeValidation}
	ErrLocked          = &Error{Code: CodeLocked}
)

func NotFound(entity string) error {
	return &Error{Code: CodeNotFound, Message: entity + " not found"}
}

func Conflict(msg string) error { return &Error{Code: CodeConflict, Message: msg} }

func Forbidden(msg string) error { return &Error{Code: CodeForbidden, Message: msg} }

func Unauthenticated(msg string) error {
	return &Error{Code: CodeUnauthenticated, Message: msg}
}

func Invariant(msg string) error { return &Error{Code: CodeInvariant, Message: msg} }

func Validation(msg string, fields map[string]string) error {
	return &Error{Code: CodeValidation, Message: msg, Fields: fields}
}

// Internal wraps an unexpected infrastructure error.
func Internal(err error) error {
	return &Error{Code: CodeInternal, Message: "internal error", Err: err}
}

// CodeOf returns the error code of err, or CodeInternal for untyped errors.
func CodeOf(err error) Code {
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return CodeInternal
}

// Validator accumulates field validation errors.
type Validator struct {
	fields map[string]string
}

// Check records msg for field when ok is false. The first failure per field wins.
func (v *Validator) Check(ok bool, field, msg string) {
	if ok {
		return
	}
	if v.fields == nil {
		v.fields = map[string]string{}
	}
	if _, exists := v.fields[field]; !exists {
		v.fields[field] = msg
	}
}

// Add records a field error unconditionally.
func (v *Validator) Add(field, msg string) { v.Check(false, field, msg) }

func (v *Validator) Valid() bool { return len(v.fields) == 0 }

// Err returns a validation error or nil.
func (v *Validator) Err() error {
	if v.Valid() {
		return nil
	}
	return &Error{Code: CodeValidation, Message: "validation failed", Fields: v.fields}
}
