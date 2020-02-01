package service

const (
	errNone = iota + 10000
	errSignedIn
	errLoginFailed
	errNoAuth
	errUnknown
)

const (
	msgSuccess     = "success"
	msgLoginFailed = "failed"
	msgNoAuth      = "no auth"
)

type loginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type httpResponse struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func (s *Service) getResponse(errCode int, msg string) httpResponse {
	return httpResponse{
		ErrorCode: errCode,
		Message:   msg,
	}
}
