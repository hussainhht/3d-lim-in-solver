//go:build js && wasm

package main

import (
	"encoding/json"
	lemin "lemin/handles"
	"syscall/js"
)

func main() {
	solve := js.FuncOf(func(_ js.Value, args []js.Value) any {
		res, err := lemin.SolveText(args[0].String())
		response := map[string]any{"result": res}
		if err != nil {
			response = map[string]any{"error": err.Error()}
		}
		data, _ := json.Marshal(response)
		return string(data)
	})
	js.Global().Set("lemInSolve", solve)
	select {}
}
