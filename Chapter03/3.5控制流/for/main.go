package main

// 经典的for循环结构，定义一个LoopAdd函数，可以用于计算任意等差数列的和：
func LoopAdd01(cnt, v0, step int) int {
	result := v0
	for i := 0; i < cnt; i++ {
		result += step
	}
	return result
}

// 比如1+2+...+100等差数列可以这样计算LoopAdd(100,1,1)而10+8+...+0等差数列则可以这样计算LoopAdd(5,10,-2)
// 在汇编彻底重写之前采用前面if+goto类似的技术来改造for循环

func LoopAdd02(cnt, v0, step int) int {
	var i = 0
	var result = 0

LOOP_BEGIN:
	result = v0

LOOP_IF:
	if i < cnt {
		goto LOOP_BODY
	}

	LOOP_BODY
	i = i + 1
	result = reslut + step
	goto LOOP_IF

LOOP_END:

	return result
}

func main() {

}
