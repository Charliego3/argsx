package argsx

import (
	"os"
	"strings"
)

// parseArgs parse os args to Value instance
func (x *argsx) parseArgs() {
	x.parser.Do(func() {
		idx := 1
		for {
			key, val := x.getKV(&idx)
			if key == "" && val == "" {
				break
			}

			ck := strings.Trim(key, "-")
			x.values[ck] = Value{key, val}
		}
	})
}

// getKV returns key value pair the key has prefix '-' value is original
func (x *argsx) getKV(idx *int) (string, string) {
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
func (x *argsx) next(idx *int) string {
	if len(os.Args)-1 < *idx {
		return ""
	}

	key := os.Args[*idx]
	*idx += 1
	return key
}
