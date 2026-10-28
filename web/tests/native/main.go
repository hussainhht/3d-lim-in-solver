// Native test oracle: serializes the same adapter used by WebAssembly.
package main

import (
	"encoding/json"
	"io"
	lemin "lemin/handles"
	"os"
)

func main() {
	text, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	res, err := lemin.SolveText(string(text))
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(res); err != nil {
		panic(err)
	}
}
