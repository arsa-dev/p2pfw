package peerconn

import (
	"io"

	"github.com/arsasm/webrtc"
)

type connWrapper struct {
	*webrtc.DataChannel
	dst io.ReadCloser
}

// NewDCConn ...
func NewDCConn(channel *webrtc.DataChannel) io.ReadWriteCloser {
	dst, src := io.Pipe()
	channel.OnMessage(func(b []byte) {
		src.Write(b)
	})
	return &connWrapper{
		DataChannel: channel,
		dst:         dst,
	}
}

func (w *connWrapper) Read(b []byte) (int, error) {
	return w.dst.Read(b)
}

func (w *connWrapper) Write(b []byte) (int, error) {
	w.DataChannel.Send(b)
	return len(b), nil
}

func (w *connWrapper) Close() error {
	if err := w.dst.Close(); err != nil {
		return err
	}
	return w.DataChannel.Close()
}
