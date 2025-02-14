package assembler

import (
	"bufio"
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"github.com/ma111e/melody/internal/engine"
	"github.com/ma111e/melody/internal/events"
	"io"
	"net/http"
)

// HTTPStreamFactory implements tcpassembly.StreamFactory
type HTTPStreamFactory struct{}

// HTTPStream will handle the actual decoding of http requests.
type HTTPStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

// New creates a new HTTPStreamFactory from the given flow data
func (h *HTTPStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &HTTPStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hstream.run() // Important... we must guarantee that data from the reader stream is read.

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &hstream.r
}

func (h *HTTPStream) run() {
	buf := bufio.NewReader(&h.r)
	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			// We must read until we see an EOF... very important!
			return
		} else if err != nil {

		} else {
			ev, _ := events.NewHTTPEvent(req, h.net, h.transport)
			engine.EventChan <- ev
		}
	}
}
