package app

import (
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/net/context"

	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/gogo/googleapis/google/rpc"
)

// HealthServer is the health server.
type HealthServer struct{}

// Check check the health.
func (s *HealthServer) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	log.Printf("handling grpc check request.")
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

// Watch is not implemented yet.
func (s *HealthServer) Watch(in *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

// AuthorizationServer is the simple authorization server.
type AuthorizationServer struct {
	AuthToken        string
	AdditionalHeader map[string]string
}

// Check check the token and set header to the response.
func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	b, err := json.MarshalIndent(req.GetAttributes().GetRequest().GetHttp().Headers, "", "  ")
	if err == nil {
		log.Printf("inbound headers: \n%s\n", string(b))
	}

	ct, err := json.MarshalIndent(req.GetAttributes().ContextExtensions, "", "  ")
	if err == nil {
		log.Printf("context extensions: \n%s\n", string(ct))
	}

	token, err := a.getAuthorizationToken(ctx, req)
	if err != nil {
		message := "Authorization header malformed or not provided."
		return responseDenied(rpc.UNAUTHENTICATED, message), nil
	}
	if token != a.AuthToken {
		message := "permission is denied."
		return responseDenied(rpc.PERMISSION_DENIED, message), nil
	}

	headers := a.getAdditionalHeaders()
	return responseOk(rpc.OK, headers), nil
}

func (a *AuthorizationServer) getAuthorizationToken(ctx context.Context, req *auth.CheckRequest) (string, error) {
	const (
		headerAuth = "authorization"
	)
	header, ok := req.GetAttributes().GetRequest().GetHttp().Headers[headerAuth]
	if !ok {
		return "", fmt.Errorf("failed to get authorization header")
	}

	return parseToken(header)
}

func (a *AuthorizationServer) getAdditionalHeaders() []*core.HeaderValueOption {
	headers := make([]*core.HeaderValueOption, 0)

	if a.AdditionalHeader == nil || len(a.AdditionalHeader) == 0 {
		return headers
	}

	for key, value := range a.AdditionalHeader {
		h := &core.HeaderValueOption{
			Header: &core.HeaderValue{
				Key:   key,
				Value: value,
			},
		}
		headers = append(headers, h)
	}
	return headers
}
