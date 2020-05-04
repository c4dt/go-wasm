package main

import (
	"errors"
	"fmt"

	"syscall/js"
)

type errorJS struct{ err error }

var _ js.Wrapper = new(errorJS)

func (e errorJS) JSValue() js.Value {
	return js.Global().Get("Error").New(e.err.Error())
}

type jsFunc func(js.Value, []js.Value) interface{}

func wrapPanic(toWrap jsFunc) jsFunc {
	return func(this js.Value, args []js.Value) (ret interface{}) {
		defer func() {
			if r := recover(); r != nil {
				ret = errorJS{fmt.Errorf("panic: %w", r)}
			}
		}()

		return toWrap(this, args)
	}
}

func increment(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return errorJS{errors.New("need one value to return")}
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
