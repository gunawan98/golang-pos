package web

type WebResponse struct {
	Success  bool        `json:"success"`
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
	Metadata interface{} `json:"metadata"`
}
