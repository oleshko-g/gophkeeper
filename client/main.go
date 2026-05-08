package main

import (
	"fmt"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/transform"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var app = struct {
	// state is used to operate the app flow
	state appState

	// cfgDir is the path to the app's data
	cfgDir string

	// config stores the config values read from the ".cfg" file stored in [app.cfgDir]
	config

	// cmd is the root command
	// during init sub commands are added based on the current [app.state]
	cmd *cobra.Command

	// is the gRPC interface to access the gophkeeper server
	*client
}{
	cmd: &cobra.Command{
		Use:   "depositor {register | authorize | connect }",
		Short: "The client app for gophkeeper",
	},
}

type appState int

const (
	cfgDirInitialized appState = iota
	cfgRead
	pubKeyRegistered
	appAuthorized
	appConnected
)

var (
	cmds = []*cobra.Command{
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
	defer func() {
		if v := recover(); v != nil {
			err, ok := v.(error)
			if !ok {
				panic(v)
			}
			cobra.CheckErr(fmt.Errorf("error during app initialization %w", err))
		}
	}()

	cobra.OnInitialize(
		initConfigDir,
		initConfig,
	)

	switch app.state {
	case cfgRead:
		app.cmd.AddCommand(register)
	}
}

var (
	register = &cobra.Command{
		Use:   "register",
		Short: "Gets the refresh token",
		Run: func(cmd *cobra.Command, args []string) {
			for input := range transform.StringFromReader(cmd.InOrStdin()) {
				fmt.Println(input)
			}
		},
	}
)

func main() {
	app.cmd.Execute()
}

func newClient(gophKeeperURL string) (*client, error) {
	conn, err := grpc.NewClient(
		gophKeeperURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &client{
		DepositorServiceClient: pb.NewDepositorServiceClient(conn),
		KeeperServiceClient:    pb.NewKeeperServiceClient(conn),
	}, nil
}

type client struct {
	pb.DepositorServiceClient
	pb.KeeperServiceClient
}
