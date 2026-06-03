package grpc

//go:generate sh -c "rm -f ./*_mock.go"
//go:generate moq -out keeper_mock.go . KeeperClient
//go:generate moq -out depositor_mock.go . DepositorClient
