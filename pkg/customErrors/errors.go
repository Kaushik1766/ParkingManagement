package customerrors

type UserNotFound struct{}

func (e UserNotFound) Error() string {
	return "user not found"
}
