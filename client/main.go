package main

import (
	"bufio"
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"log/slog"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/file"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// a is the main application struct that holds the application state and configuration.
type a struct {
	logger *slog.Logger

	// state is used to operate the app flow
	state appState

	// cfgDir is the path to the app's data
	cfgDir string

	// cfgFilePath is the path to the ".cfg" file stored in [app.cfgDir]
	cfgFilePath string

	// config stores the config values read from the cfg file stored in [app.cfgDir]
	*config

	// cmd is the root command
	// during init sub commands are added based on the current [app.state]
	cmd *cobra.Command

	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	// storager is the file-like storage object used to store deposited secrets.
	// It is initialized by [app.initDepositedSecrets]
	storager

	depositedSecrets map[string]struct{}

	// is the gRPC interface to access the gophkeeper server
	*client
}

func init() {
	godotenv.Load()

	defer func() {
		if v := recover(); v != nil {
			err, ok := v.(error)
			if !ok {
				panic(v)
			}
			cobra.CheckErr(fmt.Errorf("error during app initialization %w", err))
		}
	}()

	defer func() {
		app.logger.Info(fmt.Sprintf("app state after init is %q", app.state))
	}()

	app.initConfigDir()

	app.initConfig()

	app.initClient()

	app.initState()

	if app.state == pubKeyRegistered || app.state == appAuthorized {
		app.initKeyPair()
		app.initDepositedSecrets()
	}

	app.initState()

	switch app.state {
	case clientSetUp:
		app.cmd.AddCommand(register)
		return
	case pubKeyRegistered:
		app.cmd.AddCommand(authorize)
		return
	case appAuthorized:
		upload.SetIn(nil)
		if os.Getenv("GOPHKEEPER_DEBUG_CLIENT") == "true" {
			upload.SetIn(
				strings.NewReader("secret"),
			)
		}
		app.cmd.AddCommand(upload, list, download, delete)
		return
	}
}

var app = a{
	cmd: &cobra.Command{
		Short: "The client app for gophkeeper",
	},
	logger: slog.Default(),
}

// appState defines the state of the app.
type appState int

const (
	// cfgDirInitialized is set when the [app] configuration directory has been initialized.
	cfgDirInitialized appState = iota

	// configSet means the [config] is populated
	configSet

	// pubKeyRegistered means the [app.config.RegisteredPubKey] is populated.
	pubKeyRegistered

	// appAuthorized means the [app.config.Authorization.Token] is populated.
	appAuthorized

	// clientSetUp means the [app.client] is populated.
	clientSetUp
)

func (a appState) String() string {
	switch a {
	case cfgDirInitialized:
		return "cfgDirInitialized"
	case configSet:
		return "configSet"
	case pubKeyRegistered:
		return "pubKeyRegistered"
	case appAuthorized:
		return "appAuthorized"
	case clientSetUp:
		return "clientSetUp"

	default:
		return "unknown"
	}
}

var (
	register = &cobra.Command{
		Use:      "register",
		Short:    "Registers the client app on the gophkeeper server",
		Long:     "Generates a new RSA key pair and registers the client app on the gophkeeper server",
		RunE:     app.registerRunE,
		PostRunE: app.updateConfig,
	}
	authorize = &cobra.Command{
		Use:      "authorize",
		Short:    "Authorizes the depositor with the refresh token gotten from register command",
		RunE:     app.authorizeRunE,
		PostRunE: app.updateConfig,
	}
	upload = &cobra.Command{
		Use:      "upload",
		Short:    "Uploads a file to the gophkeeper server",
		RunE:     app.uploadRunE,
		PostRunE: app.updateDepositedSecrets,
	}
	list = &cobra.Command{
		Use:      "list",
		Short:    "Lists the data stored on the gophkeeper server",
		RunE:     app.listRunE,
		PostRunE: app.updateDepositedSecrets,
	}
	download = &cobra.Command{
		Use:   "download",
		Short: "Downloads the secret by its ID",
		RunE:  app.downloadRunE,
	}
	delete = &cobra.Command{
		Use:   "delete",
		Short: "Deletes the data stored on the gophkeeper server",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

// initConfigDir initializes the configuration directory for the [app].
func (app *a) initConfigDir() {
	userCfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	if envDevCfgDir := os.Getenv("GOPHKEEPER_DEV_CFG_DIR"); envDevCfgDir != "" {
		userCfgDir = envDevCfgDir
	}

	app.cfgDir = path.Join(userCfgDir, "gophkeeper")

	err = os.MkdirAll(app.cfgDir, dirPerm)
	if err != nil {
		panic(err)
	}

	app.logger.Info(fmt.Sprintf("config dir is initialized at %q", app.cfgDir))
}

type config struct {
	// GophKeeperURL is the URL used to set up [client]
	GophKeeperURL string `json:"keeper_url"`

	// PrivateKeyFilePath is the path to the private key file used for authentication.
	PrivateKeyFilePath string `json:"private_key_file_path"`
	PublicKeyFilePath  string `json:"public_key_file_path"`

	RegisteredPubKeyID string `json:"registered_pub_key_id"`

	Authorization struct {
		// Token is the authentication token for the [app] on the gophkeeper server.
		Token     string `json:"auth_token"`
		UpdatedAt string `json:"updated_at"`
	}
}

// initConfig initializes the configuration for the [app].
func (app *a) initConfig() {
	// read the cfg file
	const cfgFileName = ".cfg"
	app.cfgFilePath = path.Join(app.cfgDir, cfgFileName)
	cfgFile, err := os.OpenFile(app.cfgFilePath, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		panic(err)
	}
	defer cfgFile.Close()

	cfgData, err := os.ReadFile(cfgFile.Name())
	if err != nil {
		panic(err)
	}

	// init app.config
	app.config = new(config)
	if len(cfgData) == 0 {
		cfgData, err = json.MarshalIndent(config{GophKeeperURL: ":8081"}, "", "  ")
		if err != nil {
			panic(err)
		}

		_, err = cfgFile.Write(cfgData)
		if err != nil {
			panic(err)
		}
		return
	}

	err = json.Unmarshal(cfgData, app.config)
	if err != nil {
		panic(err)
	}
	app.logger.Info(fmt.Sprintf("config is initialized with values %#v", app.config))
}

func (app *a) updateConfig(cmd *cobra.Command, args []string) error {
	cfgData, err := json.MarshalIndent(app.config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(app.cfgFilePath, cfgData, filePerm)
	if err != nil {
		return err
	}
	app.logger.Info(fmt.Sprintf("config is updated to %#v", app.config))

	return nil
}

// initKeyPair initializes [app.PrivateKey] and [app.PublicKey] from the [config.PrivateKeyPath]
func (app *a) initKeyPair() {
	privKeyFileData, err := os.ReadFile(app.config.PrivateKeyFilePath)
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(privKeyFileData)
	if block == nil {
		panic(errors.New("failed to decode PEM"))
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	RSAPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		panic(errors.New("not an RSA Private Key"))
	}

	app.privateKey = RSAPrivKey
	app.publicKey = &RSAPrivKey.PublicKey

}

// initDepositedSecrets initializes [app.depositedSecrets] from the [config.DepositedKeysFile]
func (app *a) initDepositedSecrets() {
	l := app.logger.With("func", "initDepositedSecrets")

	app.depositedSecrets = make(map[string]struct{})

	const storageFileName string = "depositedSecrets.txt"
	depositedSecretsFile, err := file.RWOpenCreatePrivate(path.Join(app.cfgDir, storageFileName))
	if err != nil {
		panic(err)
	}

	app.storager = depositedSecretsFile

	s := bufio.NewScanner(app.storager)

	for c := 1; s.Scan(); c++ {
		if err := s.Err(); err != nil {
			panic(err)
		}
		line := s.Text()
		if _, err := uuid.Parse(line); err != nil {
			l.Warn("skipped an invalid deposited secret ID.", "at file line", c)
			continue
		}
		app.depositedSecrets[line] = struct{}{}
	}

	fmt.Fprintln(os.Stdout, app.depositedSecrets)
}

// updateDepositedSecrets saves [app.depositedSecrets] into [app.depositedSecretsFile]
func (app *a) updateDepositedSecrets(_ *cobra.Command, _ []string) error {
	err := app.storager.Truncate(0)
	if err != nil {
		return err
	}

	_, err = app.storager.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	for depositedSecretID, _ := range app.depositedSecrets {
		_, err = app.storager.Write([]byte(depositedSecretID + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

// initClient initializes the client for the gophkeeper server based on the [config.GophKeeperURL]
func (app *a) initClient() {
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
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
	)
	if err != nil {
		return nil, err
	}

	return &client{
		DepositorServiceClient: pb.NewDepositorServiceClient(conn),
		KeeperServiceClient:    pb.NewKeeperServiceClient(conn),
	}, nil
}

// initState initializes the state of the app based on the app's configuration and client.
func (app *a) initState() {
	if app.cfgDir != "" {
		app.state = cfgDirInitialized
	}

	if app.config == nil {
		return
	}
	app.state = configSet

	if app.client == nil {
		return
	}
	app.state = clientSetUp

	if app.config.RegisteredPubKeyID == "" {
		return
	}
	app.state = pubKeyRegistered

	if app.config.Authorization.Token == "" {
		return
	}
	app.state = appAuthorized

}

func main() {
	defer func() {
		if app.storager != nil {
			cobra.CheckErr(
				app.storager.Close(),
			)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.cmd.SetContext(ctx)
	err := app.cmd.Execute()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
