package jshttp

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

// ToResponse converts JavaScript sides Response to *http.Response.
//   - Response: https://developer.mozilla.org/docs/Web/API/Response
func ToResponse(res js.Value) (*http.Response, error) {
	status := res.Get("status").Int()
	promise := res.Call("text")
	body, err := jsutil.AwaitPromise(promise)
	if err != nil {
		return nil, err
	}
	header := ToHeader(res.Get("headers"))
	contentLength, _ := strconv.ParseInt(header.Get("Content-Length"), 10, 64)

	return &http.Response{
		Status:        strconv.Itoa(status) + " " + res.Get("statusText").String(),
		StatusCode:    status,
		Header:        header,
		Body:          io.NopCloser(strings.NewReader(body.String())),
		ContentLength: contentLength,
	}, nil
}

// ToJSResponse converts *http.Response to JavaScript sides Response class object.
func ToJSResponse(res *http.Response) js.Value {
	return newJSResponse(res.StatusCode, res.Header, res.Body, jsutil.Null)
}

// newJSResponse creates JavaScript sides Response class object.
//   - Response: https://developer.mozilla.org/docs/Web/API/Response
func newJSResponse(statusCode int, headers http.Header, body io.ReadCloser, websocket js.Value) js.Value {
	status := statusCode
	if status == 0 {
		status = http.StatusOK
	}
	respInit := jsutil.NewObject()
	respInit.Set("status", status)
	respInit.Set("statusText", http.StatusText(status))
	respInit.Set("headers", ToJSHeader(headers))
	respInit.Set("webSocket", websocket)
	if status == http.StatusSwitchingProtocols ||
		status == http.StatusNoContent ||
		status == http.StatusResetContent ||
		status == http.StatusNotModified {
		return jsutil.ResponseClass.New(jsutil.Null, respInit)
	}
	readableStream := jsutil.ConvertReaderToReadableStream(body)
	return jsutil.ResponseClass.New(readableStream, respInit)
}
