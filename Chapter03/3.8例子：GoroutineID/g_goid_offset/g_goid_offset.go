package g_goid_offset

import (
	"runtime"
	"strings"
)

var offsetDictMap = map[string]int64{
	"go1.10": 152,
	"go1.9":  152,
	"go1,8":  152,
}

var g_goid_offset = func() int64 {
	goversion := runtime.Version()
	for key, off := range offsetDictMap {
		if goversion == key || strings.HasPrefix(goversion, key) {
			return off
		}
	}
	panic("unsupport go version:" + goversion)
}()

func main() {

}
