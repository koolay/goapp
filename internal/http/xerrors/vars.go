package xerrors

var (
	// ErrDataNotExist data not found
	ErrDataNotExist = UserError{Code: 404, Message: "Notfound data"}
	// ErrInvalidParams invalid parameter
	ErrInvalidParams = UserError{Code: 400, Message: "Invalid parameters"}
	// ErrNotPermission  not permission
	ErrNotPermission = UserError{Code: 403, Message: "Not Permission"}
)
