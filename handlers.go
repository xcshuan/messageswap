package messageswap

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/libp2p/go-msgio"

	"github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pb "github.com/xcshuan/messageswap/pb"
)

// dhthandler specifies the signature of functions that handle DHT messages.
type messageHandler func(context.Context, peer.ID, *pb.Message) (*pb.Message, error)

// handleNewStream implements the network.StreamHandler
func (msnet *SwapImpl) handleNewStream(s network.Stream) {
	defer s.Reset()
	if msnet.handleNewMessage(s) {
		// Gracefully close the stream for writes.
		s.Close()
	}
}

func (msnet *SwapImpl) handleNewMessage(s network.Stream) bool {
	fmt.Println("handleNewMessage")
	ctx := msnet.context
	r := msgio.NewVarintReaderSize(s, network.MessageSizeMax)

	mPeer := s.Conn().RemotePeer()

	timer := time.AfterFunc(dhtStreamIdleTimeout, func() { s.Reset() })
	defer timer.Stop()

	for {
		var req pb.Message
		msgbytes, err := r.ReadMsg()
		if err != nil {
			fmt.Println(err)
			defer r.ReleaseMsg(msgbytes)
			if err == io.EOF {
				return true
			}
			// This string test is necessary because there isn't a single stream reset error
			// instance	in use.
			if err.Error() != "stream reset" {
				logger.Debugf("error reading message: %#v", err)
			}
			return false
		}
		err = req.XXX_Unmarshal(msgbytes)
		r.ReleaseMsg(msgbytes)
		if err != nil {
			logger.Debugf("error unmarshalling message: %#v", err)
			return false
		}

		timer.Reset(dhtStreamIdleTimeout)
		handler := msnet.handlerForMsgType(req.GetType())
		if handler == nil {
			logger.Warningf("can't handle received message of type %v", req.GetType())
			return false
		}

		resp, err := handler(ctx, mPeer, &req)
		if err != nil {
			logger.Debugf("error handling message: %v", err)
			return false
		}

		if resp == nil {
			continue
		}

		// send out response msg
		err = writeMsg(s, resp)
		if err != nil {
			logger.Debugf("error writing response: %v", err)
			return false
		}
	}
}

func (msnet *SwapImpl) handlerForMsgType(pb.Message_MessageType) messageHandler {
	return msnet.handleTest
}

func (msnet *SwapImpl) handleTest(ctx context.Context, peerID peer.ID, mes *pb.Message) (*pb.Message, error) {
	fmt.Println("Receive message from", peerID.Pretty(), "message", mes)
	return &pb.Message{
		Type: pb.Message_GET_VALUE,
		Key: pb.Message_Key{
			ID: msnet.self.Pretty(),
		},
		Value: [][]byte{[]byte("Hello, this a test")},
	}, nil
}
