package client

func isValid(status int) bool {
	switch status {
	case
		200,
		201,
		204:
		return true
	}
	return false
}
