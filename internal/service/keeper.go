package service

import pb "github.com/oleshko-g/gophkeeper/api/v1"

//go:generate moq -rm -out keeper_mock.go . Keeper
type Keeper interface {
	Namer
	pb.KeeperServiceServer
}
