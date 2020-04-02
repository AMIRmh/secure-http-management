package shm

type CustomErrors struct {
	ErrorCode int
	ErrorMsg  string
}

var (
	ErrUserNotFound = CustomErrors{ErrorCode: 1, ErrorMsg: "user not found"}
	ErrPassWrong    = CustomErrors{ErrorCode: 2, ErrorMsg: "given password is wrong"}
)
