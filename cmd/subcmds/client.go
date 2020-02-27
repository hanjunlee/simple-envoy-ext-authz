package subcmds

import (
	"context"
	"log"
	"os"

	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
)

var (
	addr    *string
	cToken   *string
)

func init() {
	addr = clientCmd.Flags().String("addr", "localhost:50051", "the address of server.")
	cToken = clientCmd.Flags().String("token", "foo", "the token of authz.")

	rootCmd.AddCommand(clientCmd)
}

var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "It help to send request to the server.",
		Run: func(cmd *cobra.Command, args []string) {
			opts := []grpc.DialOption{grpc.WithInsecure()}

			conn, err := grpc.Dial(*addr, opts...)
			if err != nil {
				os.Exit(2)
			}
			defer conn.Close()

			client := auth.NewAuthorizationClient(conn)

			res, err := client.Check(context.Background(), getReq(*cToken))
			if err != nil {
				log.Fatalf("Check has failed: %s", err)
			}

			log.Printf("the respose: %s\n", res)
		},
	}
)

func getReq(token string) *auth.CheckRequest {
	return &auth.CheckRequest{
		Attributes: &auth.AttributeContext{
			Request: &auth.AttributeContext_Request{
				Http: &auth.AttributeContext_HttpRequest{
					Headers: map[string]string{
						"authorization": "Bearer "+token,
					},
				},
			},
		},
	}
}
