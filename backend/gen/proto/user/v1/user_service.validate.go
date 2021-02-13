package user

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using https://raw.githubusercontent.com/hexdigest/gowrap/48036d11c7f254512a2e8a94b8458be52e92a899/templates/twirp_validate template

import (
	context "context"

	"github.com/twitchtv/twirp"
)

// UserServiceWithTwirpValidation implements UserService interface instrumented with arguments validation
type UserServiceWithTwirpValidation struct {
	UserService
}

// NewUserServiceWithTwirpValidation returns UserServiceWithTwirpValidation
func NewUserServiceWithTwirpValidation(base UserService) UserServiceWithTwirpValidation {
	return UserServiceWithTwirpValidation{
		UserService: base,
	}
}

// Find implements UserService
func (_d UserServiceWithTwirpValidation) Find(ctx context.Context, fp1 *FindRequest) (fp2 *FindResponse, err error) {

	if _v, _ok := interface{}(fp1).(interface{ Validate() error }); _ok {
		if err = _v.Validate(); err != nil {
			err = twirp.NewError(twirp.InvalidArgument, err.Error())
			return
		}
	}

	return _d.UserService.Find(ctx, fp1)
}
