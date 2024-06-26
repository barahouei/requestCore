package response

import "github.com/hmmftg/requestCore/webFramework"

type ResponseHandler interface {
	GetErrorsArray(message string, data any) []ErrorResponse
	HandleErrorState(err error, status int, message string, data any, w webFramework.WebFramework)
	Respond(code, status int, message string, data any, abort bool, w webFramework.WebFramework)
	RespondWithReceipt(code, status int, message string, data any, printData *Receipt, abort bool, w webFramework.WebFramework)
	OK(w webFramework.WebFramework, resp any)
	OKWithReceipt(w webFramework.WebFramework, resp any, receipt *Receipt)
	OKWithAttachment(w webFramework.WebFramework, file *FileResponse)
	Error(w webFramework.WebFramework, err ErrorState)
}

type InternalError struct {
	Desc    string
	Message any
}

type RespType int

type RespData struct {
	Code       int           `json:"code"`
	Status     int           `json:"status"`
	Message    string        `json:"message"`
	Type       RespType      `json:"type"`
	JSON       any           `json:"description"`
	PrintData  *Receipt      `json:"receipt"`
	Attachment *FileResponse `json:"attachment"`
}

const (
	Json RespType = iota
	JsonWithReceipt
	FileAttachment
)

func (e InternalError) Error() string {
	return e.Desc
}
