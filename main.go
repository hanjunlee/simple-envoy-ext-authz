package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/gogo/googleapis/google/rpc"
)

var (
	grpcport  = flag.String("grpcport", ":50051", "grpcport.")
	authToken = flag.String("token", "foo", "authenticated token.")
	conn      *grpc.ClientConn
	hs        *health.Server
)

const (
	address string = ":50051"
)

type healthServer struct{}

func (s *healthServer) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	log.Printf("Handling grpc Check request")
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *healthServer) Watch(in *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

// AuthorizationServer is the simple authorization server.
type AuthorizationServer struct {
	authToken string
}

// Check validate the token of header.
func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	log.Println(">>> Authorization called check()")

	b, err := json.MarshalIndent(req.Attributes.Request.Http.Headers, "", "  ")
	if err == nil {
		log.Println("Inbound Headers: ")
		log.Println((string(b)))
	}

	ct, err := json.MarshalIndent(req.Attributes.ContextExtensions, "", "  ")
	if err == nil {
		log.Println("Context Extensions: ")
		log.Println((string(ct)))
	}

	authHeader, ok := req.Attributes.Request.Http.Headers["authorization"]
	var splitToken []string

	if ok {
		splitToken = strings.Split(authHeader, "Bearer ")
	}
	if len(splitToken) == 2 {
		token := splitToken[1]

		if token == a.authToken {
			return &auth.CheckResponse{
				Status: &rpcstatus.Status{
					Code: int32(rpc.OK),
				},
				HttpResponse: &auth.CheckResponse_OkResponse{
					OkResponse: &auth.OkHttpResponse{
						Headers: []*core.HeaderValueOption{
							{
								Header: &core.HeaderValue{
									Key:   "x-custom-header-from-authz",
									Value: "some value",
								},
							},
						},
					},
				},
			}, nil
		} else {
			return &auth.CheckResponse{
				Status: &rpcstatus.Status{
					Code: int32(rpc.PERMISSION_DENIED),
				},
				HttpResponse: &auth.CheckResponse_DeniedResponse{
					DeniedResponse: &auth.DeniedHttpResponse{
						Status: &envoy_type.HttpStatus{
							Code: envoy_type.StatusCode_Unauthorized,
						},
						Body: "PERMISSION_DENIED",
					},
				},
			}, nil

		}

	}
	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(rpc.UNAUTHENTICATED),
		},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode_Unauthorized,
				},
				Body: "Authorization Header malformed or not provided",
			},
		},
	}, nil
}

func main() {
	flag.Parse()

	if *grpcport == "" {
		fmt.Fprintln(os.Stderr, "missing -grpcport flag (:50051)")
		flag.Usage()
		os.Exit(2)
	}

	lis, err := net.Listen("tcp", *grpcport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}
	opts = append(opts)

	s := grpc.NewServer(opts...)

	auth.RegisterAuthorizationServer(s, &AuthorizationServer{
		authToken: *authToken,
	})
	healthpb.RegisterHealthServer(s, &healthServer{})

	log.Printf("Starting gRPC Server at %s and authz token is %s", *grpcport, *authToken)
	s.Serve(lis)

}
