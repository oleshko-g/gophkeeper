package service

import (
	pb "github.com/oleshko-g/gophkeeper/api/v1"
)

type Service struct {
	Depositor
	Keeper
}

//go:generate moq -rm -out depositor_mock.go . Depositor
type Depositor = pb.DepositorServiceServer

//go:generate moq -rm -out keeper_mock.go . Keeper
type Keeper = pb.KeeperServiceServer
