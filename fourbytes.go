package fourbytes

import (
	_ "embed"
	"encoding/json"
)

// go:embed 4bytes.json
var fb []byte

var signatures map[string]map[string][]string

func init() {
	err := json.Unmarshal(fb, &signatures)
	if err != nil {
		panic(err)
	}
}

func GetCalls(sig string) []string {
	ents, ok := signatures[sig]
	if !ok {
		return nil
	}

	calls, ok := ents["calls"]
	if !ok {
		return nil
	}

	return calls
}
