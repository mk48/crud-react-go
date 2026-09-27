package util

type HttpResult struct {
	IsSuccess bool        `json:"isSuccess"`
	Error     string      `json:"error"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
}

func HttpData(result interface{}) *HttpResult {
	return &HttpResult{
		IsSuccess: true,
		Error:     "",
		Message:   "",
		Result:    result,
	}
}

func HttpError(err error, message string) *HttpResult {
	errString := ""
	if err != nil {
		errString = err.Error()
	}

	return &HttpResult{
		IsSuccess: false,
		Error:     errString,
		Message:   message,
		Result:    nil,
	}
}

func HttpErrorMessage(message string) *HttpResult {
	return &HttpResult{
		IsSuccess: false,
		Error:     "",
		Message:   message,
		Result:    nil,
	}
}

func HttpSuccessStatus() *HttpResult {
	return &HttpResult{
		IsSuccess: true,
		Error:     "",
		Message:   "",
		Result:    "",
	}
}
