package dto

type Response struct {
	Status      string `json:"status"`
	Message     string `json:"message,omitempty"`
	MessageTh   string `json:"message_th,omitempty"`
	Data        any    `json:"data,omitempty"`
	Error       string `json:"error,omitempty"`
	TotalPage   int64  `json:"total_page,omitempty"`
	TotalData   int64  `json:"total_data,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}
