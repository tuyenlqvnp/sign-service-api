package bean

const Success = "Success"
const UnexpectedError = "UnexpectedError"
const TokenInvalid = "TokenInvalid"

var CodeMessage = map[string]struct {
	Code    int
	Message string
}{
	Success:         {1, "Success"},
	UnexpectedError: {-1, "Unexpected error"},
	TokenInvalid:    {-3, "Token is invalid"},

	// 1xxx is for Borrower

	// 2xxx is for Autonomous Customer
}
