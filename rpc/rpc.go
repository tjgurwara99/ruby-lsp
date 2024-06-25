package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type rpcError string

func (err rpcError) Error() string {
	return string(err)
}

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		// TODO(taj) handle gracefully
		panic(err)
	}
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

const ErrInvalidMsg rpcError = "invalid message received from the client"

func DecodeMessage(msg []byte, v any) error {
	header, content, found := bytes.Cut(msg, []byte("\r\n\r\n"))
	if !found {
		return ErrInvalidMsg
	}
	contentLength, found := bytes.CutPrefix(header, []byte("Content-Length: "))
	if !found {
		return ErrInvalidMsg
	}
	clen, err := strconv.Atoi(string(contentLength))
	if err != nil {
		return fmt.Errorf("invalid content length in the message: %w", ErrInvalidMsg)
	}
	return json.Unmarshal(content[:clen], v)
}

type Message struct {
	Method string `json:"method"`
}

func SplitFunc(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil
	}

	contentLen := header[len("Content-Length: "):]
	clen, err := strconv.Atoi(string(contentLen))
	if err != nil {
		return 0, nil, err
	}
	if len(content) < clen {
		return 0, nil, nil
	}
	totalLen := len(header) + 4 + clen
	return totalLen, data[:totalLen], nil
}
