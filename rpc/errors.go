package rpc

type ErrorCode int

const (
	ErrParseError                  ErrorCode = -32700
	ErrInvalidRequest              ErrorCode = -32600
	ErrMethodNotFound              ErrorCode = -32601
	ErrInvalidParams               ErrorCode = -32602
	ErrInternalError               ErrorCode = -32603
	jsonrpcReservedErrorRangeStart ErrorCode = -32099
	ErrserverErrorStart            ErrorCode = jsonrpcReservedErrorRangeStart
	ErrServerNotInitialized        ErrorCode = -32002
	jsonrpcReservedErrorRangeEnd   ErrorCode = -32000
	ErrUnknownErrorCode            ErrorCode = -32001
	ErrServerErrorEnd              ErrorCode = jsonrpcReservedErrorRangeEnd
	ErrLSPReservedErrorRangeStart  ErrorCode = -32899
	ErrRequestFailed               ErrorCode = -32803
	ErrServerCancelled             ErrorCode = -32802
	ErrContentModified             ErrorCode = -32801
	ErrRequestCancelled            ErrorCode = -32800
	ErrLSPReservedErrorRangeEnd    ErrorCode = -32800
)
