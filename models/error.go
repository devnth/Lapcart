package models

type Error struct {
	Message string
}

type ResponseError struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
}

type ResponseSuccess struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}
