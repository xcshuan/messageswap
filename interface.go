package messageswap

import (
	"context"

	peer "github.com/libp2p/go-libp2p-peer"
	pb "github.com/xcshuan/messageswap/pb"
)

type Metaswap interface {
	AddAddrToPeerstore(addr string) (peer.ID, error)
	SendMessage(ctx context.Context, pmes *pb.Message, peerID peer.ID) error
	SendRequest(ctx context.Context, pmes *pb.Message, peerID peer.ID) (*pb.Message, error)
}
