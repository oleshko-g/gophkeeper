package main

import (
	"fmt"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/transform"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	cobra.OnInitialize(initConfig)
	d, err := new(a.config.GophKeeperURL)
	if err != nil {
		panic(err)
	}
	a.depositor = d
	a.cmd.AddCommand(cmds...)
}

type app struct {
	state appState
	config
	cmd *cobra.Command
	*depositor
}

type appState int

const (
	initializedCfgDir appState = iota
	readConfig
	registered
	authorized
	connected
)

var (
	a = app{
		cmd: &cobra.Command{
			Use:   "depositor {register | authorize | connect }",
			Short: "The client app for gophkeeper",
		},
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

func main() {
	a.cmd.Execute()
}

func new(gophKeeperURL string) (*depositor, error) {
	conn, err := grpc.NewClient(
		gophKeeperURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &depositor{
		DepositorServiceClient: pb.NewDepositorServiceClient(conn),
		KeeperServiceClient:    pb.NewKeeperServiceClient(conn),
	}, nil
}

type depositor struct {
	pb.DepositorServiceClient
	pb.KeeperServiceClient
}
