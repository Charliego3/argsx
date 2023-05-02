package argsx

import "sync"

type argsx struct {
	values map[string]Value
	parser sync.Once
}

var x = &argsx{
	values: make(map[string]Value),
}

// Fetch get the args value by key
//
//	os.Args = []string{"--config", "~/config/file/path.yaml"}
//	Fetch("config").String() // "~/config/file/path.yaml"
func Fetch(key string) Value {
	x.parseArgs()
	return x.values[key]
}
