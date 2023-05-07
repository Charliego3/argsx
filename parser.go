package argsx

import (
	"strings"
	"sync/atomic"
)

// parseArgs parse os args to Value instance
func (x *Argsx) parseArgs() {
	if atomic.LoadUint32(&x.done) == 1 {
		return
	}

	x.mux.Lock()
	defer x.mux.Unlock()

	if x.done == 1 {
		return
	}

	idx := 1
	for {
		key, val := x.getKV(&idx)
		if key == "" && val == "" {
			break
		}

		ck := strings.Trim(key, "-")
		x.values[ck] = Value{key, val}
	}

	atomic.StoreUint32(&x.done, 1)
}

// getKV returns key value pair the key has prefix '-' value is original
func (x *Argsx) getKV(idx *int) (string, string) {
	v := x.next(idx)
	if len(v) == 0 {
		return "", ""
	}

	if !strings.HasPrefix(v, "-") {
		return x.getKV(idx)
	}

	var key, value string
	if strings.Contains(v, "=") {
		arr := strings.SplitN(v, "=", 2)
		key = arr[0]
		value = arr[1]
	} else {
		key = v
		v = x.next(idx)
		if !strings.HasPrefix(v, "-") {
			value = v
		} else {
			*idx -= 1
		}
	}
	return key, value
}

// next get os args next value
func (x *Argsx) next(idx *int) string {
	if len(x.args)-1 < *idx {
		return ""
	}

	key := x.args[*idx]
	*idx += 1
	return key
}
