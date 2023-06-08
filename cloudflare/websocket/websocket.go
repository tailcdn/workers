package websocket

import (
	"net/http"
	"syscall/js"

	"github.com/syumai/workers/internal/jshttp"
	"github.com/syumai/workers/internal/jsutil"
)

func NewWebSocketPair() (client *WebSocket, server *WebSocket) {
	pair := jsutil.WebSocketPairClass.New()
	client = &WebSocket{Instance: pair.Index(0)}
	server = &WebSocket{Instance: pair.Index(1)}
	return
}

type WebSocket struct {
	Instance js.Value
}

func (w *WebSocket) Accept() {
	w.Instance.Call("accept")
}

func (w *WebSocket) Send(v any) {
	w.Instance.Call("send", js.ValueOf(v))
}

func (w *WebSocket) Close() {
	w.Instance.Call("close")
}

func (w *WebSocket) SetResponseWebSocket(rw http.ResponseWriter) {
	rwb := rw.(*jshttp.ResponseWriter)
	rwb.WebSocket = w.Instance
}

func (w *WebSocket) AddCloseListener(handler func(*CloseEvent)) {
	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		event := &CloseEvent{instance: args[0]}
		handler(event)
		return nil
	})
	w.Instance.Call("addEventListener", "close", cb)
}

type CloseEvent struct {
	instance js.Value
}

func (e *CloseEvent) Code() int {
	return e.instance.Get("code").Int()
}

func (e *CloseEvent) Reason() string {
	return e.instance.Get("reason").String()
}

func (e *CloseEvent) WasClean() bool {
	return e.instance.Get("wasClean").Bool()
}

func (w *WebSocket) AddMessageListener(handler func(*MessageEvent)) {
	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		event := &MessageEvent{instance: args[0]}
		handler(event)
		return nil
	})
	w.Instance.Call("addEventListener", "message", cb)
}

type MessageEvent struct {
	instance js.Value
}

func (e *MessageEvent) Data() js.Value {
	return e.instance.Get("data")
}

func (e *MessageEvent) Origin() string {
	return e.instance.Get("origin").String()
}

func (e *MessageEvent) LastEventID() string {
	return e.instance.Get("lastEventID").String()
}
