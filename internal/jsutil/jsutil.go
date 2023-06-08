package jsutil

import (
	"fmt"
	"syscall/js"
	"time"
)

var (
	Global              = js.Global()
	ObjectClass         = Global.Get("Object")
	PromiseClass        = Global.Get("Promise")
	RequestClass        = Global.Get("Request")
	ResponseClass       = Global.Get("Response")
	HeadersClass        = Global.Get("Headers")
	ArrayClass          = Global.Get("Array")
	Uint8ArrayClass     = Global.Get("Uint8Array")
	ErrorClass          = Global.Get("Error")
	ReadableStreamClass = Global.Get("ReadableStream")
	DateClass           = Global.Get("Date")
	WebSocketPairClass  = Global.Get("WebSocketPair")
	Null                = js.ValueOf(nil)
)

func NewObject() js.Value {
	return ObjectClass.New()
}

func NewUint8Array(size int) js.Value {
	return Uint8ArrayClass.New(size)
}

func NewPromise(fn js.Func) js.Value {
	return PromiseClass.New(fn)
}

// ArrayFrom calls Array.from to given argument and returns result Array.
func ArrayFrom(v js.Value) js.Value {
	return ArrayClass.Call("from", v)
}

func AwaitPromise(promiseVal js.Value) (js.Value, error) {
	resultCh := make(chan js.Value)
	errCh := make(chan error)
	var then, catch js.Func
	then = js.FuncOf(func(_ js.Value, args []js.Value) any {
		defer then.Release()
		result := args[0]
		resultCh <- result
		return js.Undefined()
	})
	catch = js.FuncOf(func(_ js.Value, args []js.Value) any {
		defer catch.Release()
		result := args[0]
		errCh <- fmt.Errorf("failed on promise: %s", result.Call("toString").String())
		return js.Undefined()
	})
	promiseVal.Call("then", then).Call("catch", catch)
	select {
	case result := <-resultCh:
		return result, nil
	case err := <-errCh:
		return js.Value{}, err
	}
}

// StrRecordToMap converts JavaScript side's Record<string, string> into map[string]string.
func StrRecordToMap(v js.Value) map[string]string {
	entries := ObjectClass.Call("entries", v)
	entriesLen := entries.Get("length").Int()
	result := make(map[string]string, entriesLen)
	for i := 0; i < entriesLen; i++ {
		entry := entries.Index(i)
		key := entry.Index(0).String()
		value := entry.Index(1).String()
		result[key] = value
	}
	return result
}

// MaybeString returns string value of given JavaScript value or returns nil if the value is undefined.
func MaybeString(v js.Value) string {
	if v.IsUndefined() {
		return ""
	}
	return v.String()
}

// MaybeDate returns time.Time value of given JavaScript Date value or returns nil if the value is undefined.
func MaybeDate(v js.Value) (time.Time, error) {
	if v.IsUndefined() {
		return time.Time{}, nil
	}
	return DateToTime(v)
}

// DateToTime converts JavaScript side's Data object into time.Time.
func DateToTime(v js.Value) (time.Time, error) {
	milli := v.Call("getTime").Float()
	return time.UnixMilli(int64(milli)), nil
}

// TimeToDate converts Go side's time.Time into Date object.
func TimeToDate(t time.Time) js.Value {
	return DateClass.New(t.UnixMilli())
}
