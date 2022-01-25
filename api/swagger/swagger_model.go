package swagger

type notAuthorized struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
