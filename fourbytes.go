package fourbytes

import (
	_ "embed"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	jsoniter "github.com/json-iterator/go"
)

const (
	// Intervall between attempts to load an ABI
	LoadInterval = 5 * time.Minute
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

type Loader interface {
	Load(string) (*abi.ABI, error)
}

type cacheEntry struct {
	abi       *abi.ABI
	lastCheck time.Time
}

type MapFunc func(*Argument) (interface{}, bool)

type Fourbytes struct {
	Loader Loader
	ArgMap MapFunc
	cache  map[string]*cacheEntry
}

func NewFourbytes() *Fourbytes {
	return &Fourbytes{
		Loader: nil,
		ArgMap: nil,
		cache:  map[string]*cacheEntry{},
	}
}

type Argument struct {
	Call  *Call
	Name  string
	Type  string
	Value interface{}
}

func (a *Argument) String() string {
	if a == nil {
		return "nil"
	}

	value := a.Value

	if f := a.Call.fb.ArgMap; f != nil {
		if v, ok := f(a); ok {
			value = v
		}
	}

	switch v := value.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v, ",")
	case common.Address:
		return strings.ToLower(v.String())
	case []common.Address:
		var list []string
		for _, addr := range v {
			list = append(list, strings.ToLower(addr.String()))
		}
		return "[" + strings.Join(list, ",") + "]"
	}

	return fmt.Sprintf("%v", a.Value)
}

type Call struct {
	fb *Fourbytes

	Address   string
	Method    string
	Arguments []*Argument
}

func (c *Call) String() string {
	if c == nil || c.Method == "" {
		return "<empty>"
	}

	var args []string
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}
	return fmt.Sprintf("%s(%s)", c.Method, strings.Join(args, ","))
}

func (fb *Fourbytes) DecodeCall(address, input string) (*Call, error) {
	address = "0x" + strings.TrimPrefix(strings.ToLower(address), "0x")
	input = strings.TrimPrefix(strings.ToLower(input), "0x")
	data, err := hex.DecodeString(input)
	if err != nil {
		return nil, err
	} else if len(data) < 4 {
		return nil, nil
	}

	ent, ok := fb.cache[address]
	if !ok {
		ent = &cacheEntry{abi: nil, lastCheck: time.Unix(0, 0)}
	}

	if ent.abi == nil && fb.Loader != nil && time.Now().After(ent.lastCheck.Add(LoadInterval)) {
		if abi, err := fb.Loader.Load(address); err == nil {
			ent.abi = abi
		}
		ent.lastCheck = time.Now()
		if !ok {
			fb.cache[address] = ent
		}
	}

	if ent.abi == nil {
		return nil, nil
	}

	method, err := ent.abi.MethodById(data[:4])
	if err != nil {
		return nil, nil
	}

	call := &Call{
		fb:      fb,
		Address: address,
		Method:  method.Name,
	}

	args, err := method.Inputs.UnpackValues(data[4:])
	if err != nil {
		return nil, err
	}

	for i, inp := range method.Inputs {
		arg := &Argument{
			Call:  call,
			Name:  inp.Name,
			Type:  inp.Type.String(),
			Value: args[i],
		}
		call.Arguments = append(call.Arguments, arg)
	}
	return call, nil
}
