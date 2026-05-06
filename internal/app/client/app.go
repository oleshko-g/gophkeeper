package client

import (
	"fmt"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/transform"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	depositor = &cobra.Command{
		Use:   "depositor {register | authorize | connect }",
		Short: "The client app for gophkeeper",
	}
	cmds = []*cobra.Command{
		&cobra.Command{
			Use:   "register",
			Short: "Gets the refresh token",
			Run: func(cmd *cobra.Command, args []string) {
				for input := range transform.StringFromReader(cmd.InOrStdin()) {
					fmt.Println(input)
				}
			},
		},
		&cobra.Command{
			Use:   "authorize",
			Short: "Authorizes the depositor with the refresh token gotten from register command",
		},
		&cobra.Command{
			Use:   "connect",
			Short: "Connect the depositor with the authentication token gooten from the authorize command",
		},
	}
)

func init() {
	depositor.AddCommand(cmds...)
	cobra.OnInitialize(initConfig)
}

func Execute() error {
	return depositor.Execute()
}

func new(cfg *Config) (*Depositor, error) {
	conn, err := grpc.NewClient(
		cfg.GophKeeperURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Depositor{
		DepositorServiceClient: pb.NewDepositorServiceClient(conn),
		KeeperServiceClient:    pb.NewKeeperServiceClient(conn),
	}, nil
}

type Depositor struct {
	pb.DepositorServiceClient
	pb.KeeperServiceClient
}

type Config struct {
	GophKeeperURL string `envconfig:"GOPHKEEPER_URL" default:":8001"`
}
