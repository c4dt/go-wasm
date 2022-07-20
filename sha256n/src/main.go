package main

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"syscall/js"
)

func toJSError(err error) js.Error {
	return js.Error{
		Value: js.Global().Get("Error").New(err.Error()),
	}
}

type sliceOfByteJS []byte

func (s sliceOfByteJS) JSValue() js.Value {
	ret := js.Global().Get("Uint8Array").New(len(s))
	js.CopyBytesToJS(ret, s)
	return ret
}

type jsFunc func(js.Value, []js.Value) interface{}

func wrapPanic(toWrap jsFunc) jsFunc {
	return func(this js.Value, args []js.Value) (ret interface{}) {
		defer func() {
			if r := recover(); r != nil {
				ret = toJSError(fmt.Errorf("panic: %w", r))
			}
		}()

		return toWrap(this, args)
	}
}

func sha256n(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return toJSError(errors.New("need the data to hash and the number of iteration"))
	}

	toHashJS := args[0]

	toHash := make([]byte, toHashJS.Length())
	js.CopyBytesToGo(toHash, toHashJS)
	iterations := args[1].Int()

	for i := 0; i < iterations; i++ {
		hashed := sha256.Sum256(toHash)
		toHash = hashed[:]
	}

	return sliceOfByteJS(toHash)
}

func main() {
	js.Global().Set("go_wasm", map[string]interface{}{
		"sha256n": js.FuncOf(wrapPanic(sha256n)),
	})

	<-make(chan struct{})
}
