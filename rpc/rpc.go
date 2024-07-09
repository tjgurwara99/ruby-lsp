package rpc

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/textproto"
	"strconv"
	"sync"
)

func Read(r *bufio.Reader) (*Request, error) {
	header, err := textproto.NewReader(r).ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	contentLength, err := strconv.ParseInt(header.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, err
	}
	var req Request
	err = json.NewDecoder(io.LimitReader(r, contentLength)).Decode(&req)
	if err != nil {
		return nil, err
	}
	if !req.IsJSONRPC() {
		return &req, ErrInvalidMsg
	}
	return &req, nil
}

type Mux struct {
	reader               *bufio.Reader
	writer               *bufio.Writer
	notificationHandlers map[string]NotificationHandler
	methodHandlers       map[string]MethodHandler
	writeLock            *sync.Mutex
	logger               *log.Logger
}

func NewMux(r io.Reader, w io.Writer, l *log.Logger) *Mux {
	reader := bufio.NewReader(r)
	writer := bufio.NewWriter(w)
	return &Mux{
		reader:               reader,
		writer:               writer,
		methodHandlers:       make(map[string]MethodHandler),
		notificationHandlers: make(map[string]NotificationHandler),
		logger:               l,
		writeLock:            &sync.Mutex{},
	}
}

func (m *Mux) HandleMethod(name string, handler MethodHandler) {
	m.methodHandlers[name] = handler
}

func (m *Mux) HandleNotification(name string, handler NotificationHandler) {
	m.notificationHandlers[name] = handler
}

func Write(w *bufio.Writer, msg Message) (err error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return
	}
	headers := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	if _, err = w.WriteString(headers); err != nil {
		return
	}
	if _, err = w.Write(body); err != nil {
		return
	}
	return w.Flush()
}

func (m *Mux) write(msg Message) error {
	m.writeLock.Lock()
	defer m.writeLock.Unlock()
	m.logger.Printf("%+v", msg)
	return Write(m.writer, msg)
}

func (m *Mux) Notify(method string, params any) error {
	n := Notification{
		Version: "2.0",
		Method:  method,
		Params:  params,
	}
	return m.write(&n)
}

func (m *Mux) Process() error {
	m.logger.Println("Started reading")
	req, err := Read(m.reader)
	if err != nil {
		m.logger.Println(err)
		return err
	}
	m.logger.Printf("method: %s", req.Method)
	m.logger.Printf("request: %s", req.Params)
	if req.IsNotification() {
		if nh, ok := m.notificationHandlers[req.Method]; ok {
			nErr := nh(req.Params)
			if nErr != nil {
				m.logger.Printf("error handling notification: %s", nErr)
			}
		}
		return nil
	}
	mh, ok := m.methodHandlers[req.Method]
	if !ok {
		m.logger.Println("method not found", req.Method)
		wErr := m.write(NewResponseError(req.ID, ErrMethodNotFound, errors.New("method not found")))
		if wErr != nil {
			m.logger.Printf("error writing to transport: %s", wErr)
		}
		return nil
	}
	result, err := mh(req.Params)
	m.logger.Printf("result created: %+v", result)
	if err != nil {
		m.logger.Printf("error happened: %s", err)
		wErr := m.write(NewResponseError(req.ID, ErrInternalError, err))
		if wErr != nil {
			m.logger.Printf("error writing to transport: %s", wErr)
		}
		return nil
	}
	m.logger.Printf("result written: %+v", result)
	wErr := m.write(NewResponse(req.ID, result))
	m.logger.Printf("result written: %+v", result)
	if wErr != nil {
		m.logger.Printf("error writing to transport: %s", wErr)
	}
	return nil
}

func NewResponse(id *json.RawMessage, result any) Message {
	return &Response{
		Version: "2.0",
		ID:      id,
		Result:  result,
	}
}

func NewResponseError(id *json.RawMessage, code ErrorCode, err error) Message {
	return &Response{
		Version: "2.0",
		Error: &Error{
			Code:    code,
			Message: err.Error(),
		},
	}
}
