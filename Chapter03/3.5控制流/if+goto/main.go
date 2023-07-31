package main

// 模拟三元表达式if函数
func If01(ok bool, a, b int) int {
	if ok {
		return a
	} else {
		return b
	}
}

// 用户汇编思维改写
func If02(ok int, a, b int) int {
	if ok == 0 {
		goto L
	}
	return a
L:
	return b
}

func main() {

}
