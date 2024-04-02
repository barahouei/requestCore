package webFramework

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type RecordData interface {
	GetId() string
	GetControlId(string) string
	GetIdList() []any
	SetId(string)
	SetValue(string)
	GetSubCategory() string
	GetValue() any
	GetValueMap() map[string]string
}

type HeaderInterface interface {
	GetId() string
	GetUser() string
	GetBranch() string
	GetBank() string
	GetPerson() string
	GetProgram() string
	GetModule() string
	GetMethod() string
	SetUser(string)
	SetBranch(string)
	SetBank(string)
	SetPerson(string)
	SetProgram(string)
	SetModule(string)
	SetMethod(string)
}

type FieldParser interface {
	Parse(string) string
}

type RequestParser interface {
	GetMethod() string
	GetPath() string
	GetHeader(target HeaderInterface) error
	GetHeaderValue(name string) string
	GetHttpHeader() http.Header
	GetBody(target any) error
	GetUri(target any) error
	GetUrlQuery(target any) error
	GetRawUrlQuery() string
	GetLocal(name string) any
	GetLocalString(name string) string
	GetUrlParam(name string) string
	GetUrlParams() map[string]string
	CheckUrlParam(name string) (string, bool)
	SetLocal(name string, value any)
	SetReqHeader(name string, value string)
	GetArgs(args ...any) map[string]string
	ParseCommand(command, title string, request RecordData, parser FieldParser) string
	SendJSONRespBody(status int, resp any) error
	Next() error
	Abort() error
	FormFile(name string) (multipart.File, *multipart.FileHeader, error)
	FormValue(name string) string
	MultiPartFile(formTagName string, handler func(multipart.File, *multipart.FileHeader) (*os.File, error)) (io.ReadCloser, error)
}

type RequestHandler interface {
	Respond(code, status int, message string, data any, abort bool)
	HandleErrorState(err error, status int, message string, data any)
}

type WebFramework struct {
	Ctx context.Context
	//Handler response.ResponseHandler
	Parser RequestParser
}
