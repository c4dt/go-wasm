package main

import (
	"errors"
	"fmt"

	"syscall/js"
)

func toJSError(err error) js.Error {
	return js.Error{
		Value: js.Global().Get("Error").New(err.Error()),
	}
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

func increment(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return toJSError(errors.New("need one value to return"))
	}

	toIncrement := args[0].Int()

	return toIncrement + 1
}

func main() {
	js.Global().Set("go_wasm", map[string]interface{}{
		"increment": js.FuncOf(wrapPanic(increment)),
	})

	<-make(chan struct{})
}
