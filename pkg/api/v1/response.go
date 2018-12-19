package v1

// Response represent a response sent to
// clients when request is successful
type Response struct {
	Data interface{} `json:"data"`
	Info string      `json:"info"`
}
