package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	mrand "math/rand"
	"os"

	"log"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/xcshuan/messageswap"

	"github.com/libp2p/go-libp2p"
	pb "github.com/xcshuan/messageswap/pb"
)

func main() {
	sourcePort := flag.Int("sp", 0, "Source port number")
	dest := flag.String("d", "", "Dest MultiAddr String")
	help := flag.Bool("help", false, "Display Help")
	debug := flag.Bool("debug", true, "Debug generated same node id on every execution.")

	flag.Parse()

	if *help {
		fmt.Printf("This program demonstrates a simple p2p chat application using libp2p\n\n")
		fmt.Printf("Usage: Run './chat -sp <SOURCE_PORT>' where <SOURCE_PORT> can be any port number. Now run './chat -d <MULTIADDR>' where <MULTIADDR> is multiaddress of previous listener host.\n")

		os.Exit(0)
	}

	var r io.Reader
	if *debug {
		// Constant random source. This will always generate the same host ID on multiple execution.
		// Don't do this in production code.
		r = mrand.New(mrand.NewSource(int64(*sourcePort)))
	} else {
		r = rand.Reader
	}

	// Creates a new RSA key pair for this host
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)

	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *sourcePort))

	h, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		panic(err)
	}
	if *dest == "" {
		// Set a function as stream handler.
		// This function  is called when a peer initiate a connection and starts a stream with this peer.
		// Only applicable on the receiving side.
		_, err := messageswap.New(h)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Run './test -d /ip4/127.0.0.1/tcp/%d/ipfs/%s (linux)' \n '.\\test.exe -d /ip4/127.0.0.1/tcp/%d/ipfs/%s (windows)' on another console. \n You can replace 127.0.0.1 with public IP as well.\n", *sourcePort, h.ID().Pretty(), *sourcePort, h.ID().Pretty())
		fmt.Printf("\nWaiting for incoming connection\n\n")
		// Hang forever
		<-make(chan struct{})
	} else {
		impl, err := messageswap.New(h)
		if err != nil {
			panic(err)
		}

		// fmt.Println("This node's multiaddress: ")
		// // IP will be 0.0.0.0 (listen on any interface) and port will be 0 (choose one for me).
		// // Although this node will not listen for any connection. It will just initiate a connect with
		// // one of its peer and use that stream to communicate.
		// fmt.Printf("%s/ipfs/%s\n", sourceMultiAddr, h.ID().Pretty())

		fmt.Println("This node's multiaddresses:")
		for _, la := range h.Addrs() {
			fmt.Printf(" - %v\n", la)
		}
		fmt.Println()

		// Turn the destination into a multiaddr.
		maddr, err := multiaddr.NewMultiaddr(*dest)
		if err != nil {
			log.Fatalln(err)
		}

		// Extract the peer ID from the multiaddr.
		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Fatalln(err)
		}

		// Add the destination's peer multiaddress in the peerstore.
		// This will be used during connection and stream creation by libp2p.
		h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

		a := &pb.Message{
			Type: pb.Message_APPEND_VALUE,
			Key: pb.Message_Key{
				ID: h.ID().Pretty(),
			},
			Value: [][]byte{[]byte("Hello")},
		}
		count := 0
		for {
			if count > 10 {
				break
			}
			rpmes, err := impl.SendRequest(context.Background(), a, info.ID)
			if err != nil {
				panic(err)
			}
			fmt.Println(rpmes)
			count++
		}
		select {}
	}
}
