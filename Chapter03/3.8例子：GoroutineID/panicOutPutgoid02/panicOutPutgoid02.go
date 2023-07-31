package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	var buf = make([]byte, 64)
	var stk = buf[:runtime.Stack(buf, false)]
	print(string(stk))
}

// 从runtime.Stack获取的字符串中就可以容易解析出goid信息
func GetGoid() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "gorotine")
	)

	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id:%v", err))
	}

	return int64(id)
}
