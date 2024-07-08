package rpc

import (
	"encoding/json"
	"errors"
)

type Message interface {
	IsJSONRPC() bool
}

type Request struct {
	Version string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params"`
}

func (m *Request) IsJSONRPC() bool {
	return m.Version == "2.0"
}

func (m *Request) IsNotification() bool {
	return m.ID == nil
}

var (
	ErrInvalidContentLengthHeader error = errors.New("invalid content length header")
	ErrInvalidMsg                 error = errors.New("invalid message")
)

type Response struct {
	Version string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Result  any              `json:"result,omitempty"`
	Error   *Error           `json:"error,omitempty"`
}

func (r *Response) IsJSONRPC() bool {
	return r.Version == "2.0"
}

type Notification struct {
	Method  string `json:"method"`
	Params  any    `json:"params"`
	Version string `json:"jsonrpc"`
}

func (n *Notification) IsJSONRPC() bool {
	return n.Version == "2.0"
}

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Data    any       `json:"data"`
}

func (e *Error) Error() string {
	return e.Message
}

type NotificationHandler func(params json.RawMessage) error

type MethodHandler func(params json.RawMessage) (result any, err error)

type InitializeParams struct {
	ClientInfo   *Info              `json:"client_info"`
	Capabilities ClientCapabilities `json:"capabilities"`
}

type Info struct {
	Name    string  `json:"name"`
	Version *string `json:"version"`
}

type ClientCapabilities struct {
}

type ServerCapabilities struct {
	TextDocumentSync TextDocumentSyncKind `json:"textDocumentSync"`
}

type TextDocumentSyncKind int

const (
	TextDocumentSyncKindNone TextDocumentSyncKind = iota
	TextDocumentSyncKindFull
	TextDocumentSyncKindIncremental
)

type InitializeResult struct {
	Capabilities ServerCapabilities
	ServerInfo   *Info
}
