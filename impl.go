package messageswap

import (
	"context"
	"fmt"
	"sync"
	"time"

	logging "github.com/ipfs/go-log"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
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
	peerstore pstore.Peerstore // Peer Registry
	smlk      sync.Mutex
	context   context.Context
}

//New 新建一个消息发送网络的功能
func New(host host.Host) (Metaswap, error) {
	if host == nil {
		return nil, ErrReadTimeout
	}
	messageNetwork := SwapImpl{
		Host:      host,
		self:      host.ID(),
		peerstore: host.Peerstore(),
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
	ipfsaddr, err := multiaddr.NewMultiaddr(addr)
	fmt.Println("AddAddrToPeerstore", ipfsaddr.String())
	if err != nil {
		return peer.ID(""), err
	}
	pid, err := ipfsaddr.ValueForProtocol(multiaddr.P_IPFS)
	fmt.Println("AddAddrToPeerstore", pid)
	if err != nil {
		return peer.ID(""), err
	}

	peerID, err := peer.IDB58Decode(pid)
	if err != nil {
		return peer.ID(""), err
	}

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := multiaddr.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerID)))
	fmt.Println("AddAddrToPeerstore", targetPeerAddr.String())
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)
	fmt.Println("AddAddrToPeerstore", targetAddr.String())
	msnet.Host.Peerstore().AddAddr(peerID, targetAddr, pstore.PermanentAddrTTL)
	return peerID, nil
}
