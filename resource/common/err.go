package err

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (appErr AppError) Error() string {
	return appErr.Message
}
