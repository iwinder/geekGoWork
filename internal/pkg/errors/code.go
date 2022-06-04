package errors

const (
	// ErrUserNotFound - 404: User not found.
	ErrNotFound int = iota + 110001

	// ErrUserAlreadyExist - 400: User already exist.
	ErrAlreadyExist
)

const (
	// ErrDatabase - 500: Database error.
	ErrDatabase int = iota + 100101
)
