package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"log/slog"

	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/transform"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// app is the main application struct that holds the app state and configuration.
var app = struct {
	logger *slog.Logger

	// state is used to operate the app flow
	state appState

	// cfgDir is the path to the app's data
	cfgDir string

	// config stores the config values read from the ".cfg" file stored in [app.cfgDir]
	*config

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
	config: &config{},
	logger: slog.Default(),
}

// appState defines the state of the app.
type appState int

const (
	// cfgDirInitialized is set when the [app] configuration directory has been initialized.
	cfgDirInitialized appState = iota

	// configSet means the [app.config] is populated
	configSet

	// pubKeyRegistered means the [app.config.RegisteredPubKey] is populated.
	pubKeyRegistered

	// appAuthorized means the [app.config.AuthToken] is populated.
	appAuthorized

	// appConnected means the app has been connected to the server.
	appConnected
)

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
	authorize = &cobra.Command{
		Use:   "authorize",
		Short: "Authorizes the depositor with the refresh token gotten from register command",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	connect = &cobra.Command{
		Use:   "connect",
		Short: "Connect the depositor with the authentication token gotten from the authorize command",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	upload = &cobra.Command{
		Use:   "upload",
		Short: "Uploads a file to the gophkeeper server",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	list = &cobra.Command{
		Use:   "list",
		Short: "Lists the data stored on the gophkeeper server",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	delete = &cobra.Command{
		Use:   "delete",
		Short: "Deletes the data stored on the gophkeeper server",
		Run: func(cmd *cobra.Command, args []string) {
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

	initConfigDir()

	initConfig()

	initClient()

	switch app.state {
	case configSet:
		app.cmd.AddCommand(register)
		app.cmd.Use = "depositor {register}"
		return
	case pubKeyRegistered:
		app.cmd.AddCommand(authorize)
		app.cmd.Use = "depositor {authorize}"
		return
	case appAuthorized:
		app.cmd.AddCommand(connect)
		app.cmd.Use = "depositor {connect}"
		return
	case appConnected:
		app.cmd.AddCommand(
			list,
			upload,
			delete,
		)
		app.cmd.Use = "depositor {upload | list | delete}"
	}
}

// initConfigDir initializes the configuration directory for the [app].
func initConfigDir() {
	userCfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	app.cfgDir = path.Join(userCfgDir, "gophkeeper")

	err = os.MkdirAll(app.cfgDir, dirPerm)
	if err != nil {
		panic(err)
	}

	app.state = cfgDirInitialized
	app.logger.Info(fmt.Sprintf("config dir is initialized at %q", app.cfgDir))
}

type config struct {
	GophKeeperURL string `json:"keeper_url"`

	// RegisteredPubKey is the public key registered on the gophkeeper server.
	// The key is the public key, the value is the refresh token with which the client must [authorize] the [app].
	RegisteredPubKey map[string]string `json:"registered_pub_key"`

	// AuthToken is the authentication token for the [app] on the gophkeeper server.
	// The key is the refresh token, the value is the authentication token.
	AuthToken map[string]string `json:"auth_token"`

	// Session is the session token for the [app] on the gophkeeper server.
	Session pb.Session `json:"session"`
}

// initConfig initializes the configuration for the [app].
func initConfig() {
	cfgFile, err := os.OpenFile(path.Join(app.cfgDir, ".cfg"), os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		panic(err)
	}

	cfgData, err := os.ReadFile(cfgFile.Name())
	if err != nil {
		panic(err)
	}

	if len(cfgData) == 0 {
		cfgData, err = json.Marshal(config{
			GophKeeperURL: ":8080",
		})
		if err != nil {
			panic(err)
		}

		_, err = cfgFile.Write(cfgData)
		if err != nil {
			panic(err)
		}
		app.state = configSet
	}

	err = json.Unmarshal(cfgData, app.config)
	if err != nil {
		panic(err)
	}
	app.logger.Info(fmt.Sprintf("config is initialized with values %#v", app.config))
	app.state = configSet
}

// initClient initializes the client for the gophkeeper server.
func initClient() {
	client, err := newClient(app.config.GophKeeperURL)
	if err != nil {
		panic(err)
	}
	app.client = client
}

type client struct {
	pb.DepositorServiceClient
	pb.KeeperServiceClient
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

func main() {
	err := app.cmd.Execute()
	if err != nil {
		cobra.CheckErr(err)
	}
}
