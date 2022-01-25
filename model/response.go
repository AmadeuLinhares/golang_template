package model

// Response is a default model
type Response struct {
	Data interface{}            `json:"data,omitempty" swaggerignore:"true"`
	Meta map[string]interface{} `json:"metadata,omitempty" swaggerignore:"true"`
	Err  error                  `json:"error,omitempty" swaggerignore:"true"`
}
