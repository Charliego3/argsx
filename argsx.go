package argsx

import (
	"os"
	"sync"
	"sync/atomic"
)

type Argsx struct {
	args   []string
	values map[string]Value
	done   uint32
	mux    sync.Mutex
}

// New returns arg parser with os.Args
func New() *Argsx {
	return NewWithArgs(os.Args)
}

// NewWithArgs returns arg parser with custom args
func NewWithArgs(args []string) *Argsx {
	return &Argsx{
		args:   args,
		values: make(map[string]Value),
	}
}

// Fetch get the args value by key
func (x *Argsx) Fetch(key string) Value {
	x.parseArgs()
	return x.values[key]
}

// SetArgs replace the old args
func (x *Argsx) SetArgs(args []string) {
	x.args = args
	atomic.StoreUint32(&x.done, 0)
}

var dx = New()

// SetArgs replace old args
func SetArgs(args []string) {
	dx.SetArgs(args)
}

// Fetch get the args value by key
//
//	os.Args = []string{"--config", "~/config/file/path.yaml"}
//	Fetch("config").String() // "~/config/file/path.yaml"
func Fetch(key string) Value {
	return dx.Fetch(key)
}
