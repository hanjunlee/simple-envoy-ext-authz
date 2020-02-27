package cmd

import (
	"log"
	"net"
	"strings"

	"github.com/hanjunlee/simple-envoy-ext-authz/internal/app"

	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/spf13/cobra"
)

var (
	port    *string
	token   *string
	headers *[]string
)

func init() {
	port = serverCmd.Flags().String("port", ":50051", "the port of server.")
	token = serverCmd.Flags().String("token", "foo", "the valid authorized token.")
	headers = serverCmd.Flags().StringArrayP("header", "H", nil, "add headers into the response, e.g foo=bar.")

	rootCmd.AddCommand(serverCmd)
}

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run the external auth server of Envoy.",
		Run: func(cmd *cobra.Command, args []string) {
			opts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}

			s := grpc.NewServer(opts...)
			auth.RegisterAuthorizationServer(s, &app.AuthorizationServer{
				AuthToken:        *token,
				AdditionalHeader: parseHeaders(*headers),
			})
			healthpb.RegisterHealthServer(s, &app.HealthServer{})

			log.Printf("Starting gRPC Server at %s, valid token: %s, headers: %v.", *port, *token, parseHeaders(*headers))
			lis, _ := net.Listen("tcp", *port)
			s.Serve(lis)
		},
	}
)

func parseHeaders(headers []string) map[string]string {
	const (
		sep = "="
	)

	ret := make(map[string]string)
	for _, header := range headers {
		ss := strings.Split(header, sep)
		key, value := ss[0], ss[1]
		ret[key] = value
	}
	return ret
}
