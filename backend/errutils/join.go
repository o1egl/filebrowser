package errutils

import "errors"

func JoinFn(err *error, fn func() error) {
	gotErr := fn()
	switch *err {
	case nil:
		*err = gotErr
	default:
		*err = errors.Join(*err, gotErr)
	}
}
