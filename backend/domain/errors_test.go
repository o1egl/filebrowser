package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundError_Error(t *testing.T) {
	testCases := map[string]struct {
		notFoundError NotFoundError
		want          string
	}{
		"with one constraint": {
			notFoundError: NotFoundError{
				cause:           errors.New("causer"),
				resource:        ResourceUser,
				constraintPairs: []string{"id", "123"},
			},
			want: "user with id=123 not found",
		},
		"with multiple constraints": {
			notFoundError: NotFoundError{
				cause:           errors.New("causer"),
				resource:        ResourceUser,
				constraintPairs: []string{"id", "123", "tenant", "foo"},
			},
			want: "user with id=123, tenant=foo not found",
		},
		"with odd constraints count": {
			notFoundError: NotFoundError{
				cause:           errors.New("causer"),
				resource:        ResourceUser,
				constraintPairs: []string{"id", "123", "tenant"},
			},
			want: "user with id=123, tenant= not found",
		},
		"with no constraints": {
			notFoundError: NotFoundError{
				cause:    errors.New("causer"),
				resource: ResourceUser,
			},
			want: "user not found",
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tt.notFoundError.Error()
			assert.Equal(t, tt.want, got)
		})
	}
}
