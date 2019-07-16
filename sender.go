package messageswap

import (
	"bufio"
	"context"
	"io"
	"sync"
	"time"

	ggio "github.com/gogo/protobuf/io"
	"github.com/libp2p/go-libp2p-core/helpers"
	host "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-msgio"
	pb "github.com/xcshuan/messageswap/pb"
)

var dhtReadMessageTimeout = time.Minute
var dhtStreamIdleTimeout = 10 * time.Minute

var writerPool = sync.Pool{
	New: func() interface{} {
		w := bufio.NewWriter(nil)
		return &bufferedDelimitedWriter{
			Writer:      w,
			WriteCloser: ggio.NewDelimitedWriter(w),
		}
	},
}

func writeMsg(w io.Writer, mes *pb.Message) error {
	bw := writerPool.Get().(*bufferedDelimitedWriter)
	bw.Reset(w)
	err := bw.WriteMsg(mes)
	if err == nil {
		err = bw.Flush()
	}
	bw.Reset(nil)
	writerPool.Put(bw)
	return err
}

type bufferedDelimitedWriter struct {
	*bufio.Writer
	ggio.WriteCloser
}

type bufferedWriteCloser interface {
	ggio.WriteCloser
	Flush() error
}

type messageSender struct {
	s        network.Stream
	r        msgio.ReadCloser
	lk       sync.Mutex
	p        peer.ID
	host     host.Host
	metaswap *SwapImpl
	// invalid   bool
	singleMes int
}

func (msnet *SwapImpl) NewMessageSender(ctx context.Context, p peer.ID) (*messageSender, error) {
	return &messageSender{
		host:     msnet.Host,
		p:        p,
		metaswap: msnet,
	}, nil
}

func (s *messageSender) Close() error {
	return helpers.FullClose(s.s)
}

func (s *messageSender) Reset() error {
	return s.s.Reset()
}

func (msnet *SwapImpl) newStreamToPeer(ctx context.Context, p peer.ID) (network.Stream, error) {
	return msnet.Host.NewStream(ctx, p, protocolMessageswapOne)
}

func newBufferedDelimitedWriter(str io.Writer) bufferedWriteCloser {
	w := bufio.NewWriter(str)
	return &bufferedDelimitedWriter{
		Writer:      w,
		WriteCloser: ggio.NewDelimitedWriter(w),
	}
}

func (w *bufferedDelimitedWriter) Flush() error {
	return w.Writer.Flush()
}

func (ms *messageSender) SendMessage(ctx context.Context, pmes *pb.Message) error {
	ms.lk.Lock()
	defer ms.lk.Unlock()
	retry := false
	for {
		if ms.s != nil {
			return nil
		}
		nstr, err := ms.host.NewStream(ctx, ms.p, protocolMessageswapOne)
		if err != nil {
			return err
		}

		ms.r = msgio.NewVarintReaderSize(nstr, network.MessageSizeMax)
		ms.s = nstr

		if err := ms.writeMsg(pmes); err != nil {
			ms.s.Reset()
			ms.s = nil

			if retry {
				logger.Info("error writing message, bailing: ", err)
				return err
			}
			logger.Info("error writing message, trying again: ", err)
			retry = true
			continue
		}

		logger.Event(ctx, "dhtSentMessage", ms.metaswap.self, ms.p, pmes)

		if ms.singleMes > streamReuseTries {
			go helpers.FullClose(ms.s)
			ms.s = nil
		} else if retry {
			ms.singleMes++
		}

		return nil
	}
}

func (ms *messageSender) SendRequest(ctx context.Context, pmes *pb.Message) (*pb.Message, error) {
	ms.lk.Lock()
	defer ms.lk.Unlock()
	retry := false
	for {
		if ms.s != nil {
			return nil, nil
		}
		nstr, err := ms.host.NewStream(ctx, ms.p, protocolMessageswapOne)
		if err != nil {
			return nil, err
		}

		ms.r = msgio.NewVarintReaderSize(nstr, network.MessageSizeMax)
		ms.s = nstr

		if err := ms.writeMsg(pmes); err != nil {
			ms.s.Reset()
			ms.s = nil

			if retry {
				logger.Info("error writing message, bailing: ", err)
				return nil, err
			}
			logger.Info("error writing message, trying again: ", err)
			retry = true
			continue
		}

		mes := new(pb.Message)
		if err := ms.ctxReadMsg(ctx, mes); err != nil {
			ms.s.Reset()
			ms.s = nil

			if retry {
				logger.Info("error reading message, bailing: ", err)
				return nil, err
			}
			logger.Info("error reading message, trying again: ", err)
			retry = true
			continue
		}

		logger.Event(ctx, "dhtSentMessage", ms.metaswap.self, ms.p, pmes)

		if ms.singleMes > streamReuseTries {
			go helpers.FullClose(ms.s)
			ms.s = nil
		} else if retry {
			ms.singleMes++
		}

		return mes, nil
	}
}

func (ms *messageSender) writeMsg(pmes *pb.Message) error {
	return writeMsg(ms.s, pmes)
}

func (ms *messageSender) ctxReadMsg(ctx context.Context, mes *pb.Message) error {
	errc := make(chan error, 1)
	go func(r msgio.ReadCloser) {
		bytes, err := r.ReadMsg()
		defer r.ReleaseMsg(bytes)
		if err != nil {
			errc <- err
			return
		}
		errc <- mes.XXX_Unmarshal(bytes)
	}(ms.r)

	t := time.NewTimer(dhtReadMessageTimeout)
	defer t.Stop()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return ErrReadTimeout
	}
}
