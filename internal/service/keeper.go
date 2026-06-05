package service

import pb "github.com/oleshko-g/gophkeeper/api/v1"

type Keeper interface {
	Namer
	pb.KeeperServiceServer
}
