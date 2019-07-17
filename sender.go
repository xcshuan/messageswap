package messageswap

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	ggio "github.com/gogo/protobuf/io"
	"github.com/libp2p/go-libp2p-core/helpers"
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

type messageSender struct {
	s         network.Stream
	r         msgio.ReadCloser
	lk        sync.Mutex
	p         peer.ID
	metaswap  *SwapImpl
	invalid   bool
	singleMes int
}

func (msnet *SwapImpl) NewMessageSender(ctx context.Context, p peer.ID) (*messageSender, error) {
	msnet.smlk.Lock()
	ms, ok := msnet.strmap[p]
	if ok {
		msnet.smlk.Unlock()
		return ms, nil
	}
	ms = &messageSender{p: p, metaswap: msnet}
	msnet.strmap[p] = ms
	msnet.smlk.Unlock()

	if err := ms.prepOrInvalidate(ctx); err != nil {
		msnet.smlk.Lock()
		defer msnet.smlk.Unlock()

		if msCur, ok := msnet.strmap[p]; ok {
			// Changed. Use the new one, old one is invalid and
			// not in the map so we can just throw it away.
			if ms != msCur {
				return msCur, nil
			}
			// Not changed, remove the now invalid stream from the
			// map.
			delete(msnet.strmap, p)
		}
		// Invalid but not in map. Must have been removed by a disconnect.
		return nil, err
	}
	// All ready to go.
	return ms, nil
}

func (ms *messageSender) Close() error {
	return helpers.FullClose(ms.s)
}

func (ms *messageSender) Reset() error {
	return ms.s.Reset()
}

func (ms *messageSender) prep(ctx context.Context) error {
	if ms.invalid {
		return fmt.Errorf("message sender has been invalidated")
	}
	if ms.s != nil {
		return nil
	}

	nstr, err := ms.metaswap.Host.NewStream(ctx, ms.p, protocolMessageswapOne)
	if err != nil {
		return err
	}

	ms.r = msgio.NewVarintReaderSize(nstr, network.MessageSizeMax)
	ms.s = nstr

	return nil
}

func (ms *messageSender) prepOrInvalidate(ctx context.Context) error {
	ms.lk.Lock()
	defer ms.lk.Unlock()
	if err := ms.prep(ctx); err != nil {
		ms.invalidate()
		return err
	}
	return nil
}

// invalidate is called before this messageSender is removed from the strmap.
// It prevents the messageSender from being reused/reinitialized and then
// forgotten (leaving the stream open).
func (ms *messageSender) invalidate() {
	ms.invalid = true
	if ms.s != nil {
		ms.s.Reset()
		ms.s = nil
	}
}

func (msnet *SwapImpl) newStreamToPeer(ctx context.Context, p peer.ID) (network.Stream, error) {
	return msnet.Host.NewStream(ctx, p, protocolMessageswapOne)
}

func (w *bufferedDelimitedWriter) Flush() error {
	return w.Writer.Flush()
}

func (ms *messageSender) SendMessage(ctx context.Context, pmes *pb.Message) error {
	ms.lk.Lock()
	defer ms.lk.Unlock()
	retry := false
	for {
		if err := ms.prep(ctx); err != nil {
			return err
		}

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
		if err := ms.prep(ctx); err != nil {
			return nil, err
		}

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
