package libCallApi

import (
	"time"

	"github.com/hmmftg/requestCore/response"
)

type CallParam *CallParamData
type CallParamData struct {
	Parameters  map[string]any
	Headers     map[string]string
	Api         RemoteApi
	Timeout     time.Duration
	Method      string
	Path        string
	Query       string
	QueryStack  *[]string
	ValidateTls bool
	EnableLog   bool
	JsonBody    any
}

type CallResult[RespType any] struct {
	Resp   *RespType
	WsResp *response.WsRemoteResponse
	Status *CallResp
	Error  response.ErrorState
}

func Call[RespType any](param CallParam) CallResult[RespType] {
	if len(*param.QueryStack) > 0 {
		param.Query = (*param.QueryStack)[0]
		if len(*param.QueryStack) > 1 {
			*param.QueryStack = (*param.QueryStack)[1:]
		} else {
			*param.QueryStack = nil
		}
	}
	callData := CallData{
		Api:       param.Api,
		Path:      param.Path + param.Query,
		Method:    param.Method,
		Headers:   param.Headers,
		SslVerify: !param.ValidateTls,
		EnableLog: param.EnableLog,
		Timeout:   param.Timeout,
		Req:       param.JsonBody,
	}
	resp, wsResp, callResp, err := ConsumeRest[RespType](callData)
	return CallResult[RespType]{resp, wsResp, callResp, err}
}