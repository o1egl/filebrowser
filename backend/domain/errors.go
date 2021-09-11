//go:generate go-enum --sql --marshal --nocase --names --file $GOFILE
package domain

import (
	"fmt"
	"strings"
)

/*
ENUM(
user
file
)
*/
type Resource int

type NotFoundError struct {
	cause           error
	resource        Resource
	constraintPairs []string
}

func NewNotFoundError(cause error, resource Resource, constraintPairs ...string) *NotFoundError {
	return &NotFoundError{cause: cause, resource: resource, constraintPairs: constraintPairs}
}

func (n NotFoundError) Error() string {
	var constraint string
	if len(n.constraintPairs) > 0 {
		sb := strings.Builder{}
		sb.WriteString(" with ")
		for i, part := range n.constraintPairs {
			isKey := i%2 == 0
			sb.WriteString(part)
			if isKey {
				sb.WriteByte('=')
			}
			if !isKey && i != len(n.constraintPairs)-1 {
				sb.WriteString(", ")
			}
		}
		constraint = sb.String()
	}
	return fmt.Sprintf("%s%s not found", n.resource.String(), constraint)
}

func (n NotFoundError) Unwrap() error {
	return n.cause
}

type AccessDeniedError struct {
	subject string
	object  string
}

func NewAccessDeniedError(subject string, object string) *AccessDeniedError {
	return &AccessDeniedError{subject: subject, object: object}
}

func (a AccessDeniedError) Error() string {
	return fmt.Sprintf("%s has no right to access %s", a.subject, a.object)
}
