package fourbytes

import (
	_ "embed"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

//go:embed 4bytes.json
var fb []byte

var signatures map[string]map[string][]string

func init() {
	err := json.Unmarshal(fb, &signatures)
	if err != nil {
		panic(err)
	}
}

func JSON() []byte {
	return fb
}

func Signatures() map[string]map[string][]string {
	return signatures
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
