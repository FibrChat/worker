package request

type ResponseCode int

const (
	CodeSuccess ResponseCode = iota
	CodeInternal
)

type Response struct {
	Code  ResponseCode `json:"code"`
	Error string       `json:"error,omitempty"`
}

type UsersResponse struct {
	Users []string `json:"users"`
}
