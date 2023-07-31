package main

import (
	"fmt"
	"reflect"
	"sync"
)

// 有了goid之后，构造goroutine局部存储非常容易，先定义一个gls包提供goid特性：
var gls struct {
	m map[int64]map[interface{}]interface{}
	sync.Mutex
}

func init() {
	gls.m = make(map[int64]map[interface{}]interface{})
}

// 基于g返回的接口，就可以容易获取goid
func GetGoid() int64 {
	g := getg()
	gid := reflect.ValueOf(g).FieldByName("goid").Int()
	return goid
}

// gls包变量简单包赚了map，同时通过sync.Mutex互斥量支持并发访问
// 然后定义一个getMap内部函数，用于获取每个Goroutine字节的map
func getMap() map[interface{}]interface{} {
	gls.Lock()
	defer gls.Unlock()

	goid := GetGoid()
	if m, _ := gls.m[goid]; m != nil {
		return m
	}

	m := make(map[interface{}]interface{})
	gls.m[goid] = m
	return m
}

// 获取到goroutine私有的map之后，就是正常的增、删、改操作接口：
func Get(key interface{}) interface{} {
	return getMap()[key]
}

func Put(key interface{}, v interface{}) {
	getMap()[key] = v
}

func Delete(key interface{}) {
	delete(getMap(), key)
}

// 最后提供一个Clean函数，用于释放goroutine对应的map资源
func Clean() {
	gls.Lock()
	defer gls.Unlock()

	delete(gls.m, GetGoid())
}

// 使用局部存储简单的例子
func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			defer gls.Clean()

			defer func() {
				fmt.Printf("%d: number = %d\n", idx, gls.Get("number"))
			}()
			gls.Put("number", idx+100)
		}(i)
	}
	wg.Wait()
}
