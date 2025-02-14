package httpparser

import (
	"bytes"
	"fmt"
	"github.com/ma111e/melody/internal/config"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/c2h5oh/datasize"
)

// GetBodyPayload extract the body of an http.Request without striping it
func GetBodyPayload(r *http.Request) ([]byte, error) {
	var buf bytes.Buffer
	var b bytes.Buffer
	var dest io.Writer = &b

	if r.Body != nil {
		if r.Header.Get("Content-Length") == "" {
			return []byte{}, nil
		}
		iContentLength, err := strconv.ParseUint(r.Header.Get("Content-Length"), 10, 64)

		if err != nil {
			return []byte{}, fmt.Errorf("request data not logged (failed to parse Content-Length as uint64 : %s)", err)
		}

		if iContentLength > config.Cfg.MaxPOSTDataSize {
			return []byte{}, fmt.Errorf("request data not logged (over %s : %s)", datasize.ByteSize(config.Cfg.MaxPOSTDataSize).HumanReadable(), (datasize.ByteSize(iContentLength) * datasize.B).HumanReadable())
		}

		if _, err := buf.ReadFrom(r.Body); err != nil {
			return []byte{}, fmt.Errorf("failed to parse request body : got error when reading from body [%s]", err.Error())
		}

		if err := r.Body.Close(); err != nil {
			// Send read body in such case
			return buf.Bytes(), fmt.Errorf("failed to parse request body : got error while closing body [%s]", err.Error())
		}

		bodyReader := ioutil.NopCloser(bytes.NewReader(buf.Bytes()))

		chunked := len(r.TransferEncoding) > 0 && r.TransferEncoding[0] == "chunked"

		if chunked {
			dest = httputil.NewChunkedWriter(dest)
		}

		if _, err := io.Copy(dest, bodyReader); err != nil {
			// Send read body in such case
			return buf.Bytes(), fmt.Errorf("failed to parse request body: got error while copying the read body [%s]", err.Error())
		}

		if chunked {
			_ = dest.(io.Closer).Close()
			_, _ = io.WriteString(&b, "\r\n")
		}
	}

	return b.Bytes(), nil
}
