package app

import (
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/gogo/googleapis/google/rpc"
)

func responseOk(code rpc.Code, headers []*core.HeaderValueOption) *auth.CheckResponse {
	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(code),
		},
		HttpResponse: &auth.CheckResponse_OkResponse{
			OkResponse: &auth.OkHttpResponse{
				Headers: headers,
			},
		},
	}
}

func responseDenied(code rpc.Code, message string) *auth.CheckResponse {
	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(code),
		},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode_Unauthorized,
				},
				Body: message,
			},
		},
	}
}
