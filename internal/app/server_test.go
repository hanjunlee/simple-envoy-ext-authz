package app

import (
	"context"
	"reflect"
	"testing"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/gogo/googleapis/google/rpc"
)

func getReq(headers map[string]string) *auth.CheckRequest {
	return &auth.CheckRequest{
		Attributes: &auth.AttributeContext{
			Request: &auth.AttributeContext_Request{
				Http: &auth.AttributeContext_HttpRequest{
					Headers: headers,
				},
			},
		},
	}
}

func TestAuthorizationServer_Check(t *testing.T) {
	type fields struct {
		AuthToken        string
		AdditionalHeader map[string]string
	}
	type args struct {
		ctx context.Context
		req *auth.CheckRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *auth.CheckResponse
		wantErr bool
	}{
		{
			name: "bad token",
			fields: fields{
				AuthToken:        "foo",
				AdditionalHeader: map[string]string{},
			},
			args: args{
				ctx: context.Background(),
				req: getReq(map[string]string{
					"authorization": "Bearer bar",
				}),
			},
			want: responseDenied(rpc.PERMISSION_DENIED, "permission is denied."),
		},
		{
			name: "authorized token",
			fields: fields{
				AuthToken:        "foo",
				AdditionalHeader: map[string]string{},
			},
			args: args{
				ctx: context.Background(),
				req: getReq(map[string]string{
					"authorization": "Bearer foo",
				}),
			},
			want: responseOk(rpc.OK, []*core.HeaderValueOption{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationServer{
				AuthToken:        tt.fields.AuthToken,
				AdditionalHeader: tt.fields.AdditionalHeader,
			}
			got, err := a.Check(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthorizationServer.Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthorizationServer.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorizationServer_getAuthorizationToken(t *testing.T) {
	type fields struct {
		AuthToken        string
		AdditionalHeader map[string]string
	}
	type args struct {
		ctx context.Context
		req *auth.CheckRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "bad token",
			fields: fields{
				AuthToken: "foo",
			},
			args: args{
				ctx: context.Background(),
				req: getReq(map[string]string{
					"authorization": "foo bar",
				}),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "authorized token",
			fields: fields{
				AuthToken: "foo",
			},
			args: args{
				ctx: context.Background(),
				req: getReq(map[string]string{
					"authorization": "Bearer foo",
				}),
			},
			want:    "foo",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationServer{
				AuthToken:        tt.fields.AuthToken,
				AdditionalHeader: tt.fields.AdditionalHeader,
			}
			got, err := a.getAuthorizationToken(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthorizationServer.getAuthorizationToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthorizationServer.getAuthorizationToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
