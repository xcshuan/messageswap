package messageswap

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	logging "github.com/ipfs/go-log"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	host "github.com/libp2p/go-libp2p-host"
	protocol "github.com/libp2p/go-libp2p-protocol"
	"github.com/multiformats/go-multiaddr"
	pb "github.com/xcshuan/messageswap/pb"
)

var protocolMessageswapOne protocol.ID = "/ipfs/messageswap/1.0.0"

// streamReuseTries is the number of times we will try to reuse a stream to a
// given peer before giving up and reverting to the old one-message-per-stream
// behaviour.
const streamReuseTries = 3

var ErrReadTimeout = fmt.Errorf("timed out reading response")

var logger = logging.Logger("Messageswap_network")

var sendMessageTimeout = time.Minute * 10

type SwapImpl struct {
	self      peer.ID
	Host      host.Host
	peerstore peerstore.Peerstore // Peer Registry

	smlk    sync.Mutex
	strmap  map[peer.ID]*messageSender //用于stream复用
	context context.Context
}

//New 新建一个消息发送网络的功能
func New(host host.Host) (Metaswap, error) {
	if host == nil {
		return nil, errors.New("Host is nil")
	}
	messageNetwork := SwapImpl{
		Host:      host,
		self:      host.ID(),
		peerstore: host.Peerstore(),
		strmap:    make(map[peer.ID]*messageSender),
	}
	host.SetStreamHandler(protocolMessageswapOne, messageNetwork.handleNewStream)
	return &messageNetwork, nil
}

func (msnet *SwapImpl) SendMessage(ctx context.Context, pmes *pb.Message, peerID peer.ID) error {
	sender, err := msnet.NewMessageSender(ctx, peerID)
	if err != nil {
		return err
	}
	return sender.SendMessage(ctx, pmes)
}

func (msnet *SwapImpl) SendRequest(ctx context.Context, pmes *pb.Message, peerID peer.ID) (*pb.Message, error) {
	sender, err := msnet.NewMessageSender(ctx, peerID)
	if err != nil {
		return nil, err
	}
	return sender.SendRequest(ctx, pmes)
}

func (msnet *SwapImpl) AddAddrToPeerstore(addr string) (peer.ID, error) {
	// Turn the destination into a multiaddr.
	maddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return peer.ID(""), err
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return peer.ID(""), err
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	msnet.Host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	return info.ID, nil
}
