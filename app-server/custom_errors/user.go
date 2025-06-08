package custom_errors

import "errors"

var ErrUserIdRequired = errors.New("user ID is required")
