package handlers

import (
	"encoding/json"

	"github.com/tjgurwara99/ruby-lsp/rpc"
)

func Initialize(params json.RawMessage) (any, error) {
	var initializeParams rpc.InitializeParams
	if err := json.Unmarshal(params, &initializeParams); err != nil {
		return nil, err
	}
	result := rpc.InitializeResult{
		Capabilities: rpc.ServerCapabilities{
			TextDocumentSync: rpc.TextDocumentSyncKindFull,
		},
		ServerInfo: &rpc.Info{
			Name: "ruby-lsp",
		},
	}
	return result, nil
}
