package main

import "sync"

// 计划通过rpc构造一个简单的内存KV数据库:
type KVStoreService struct {
	m      map[string]string
	filter map[string]func(key string)
	mu     sync.Mutex
}

func NewKVStoreService() *KVStoreService {
	return &KVStoreService{
		m:      make(map[string]string),
		filter: make(map[string]func(key string)),
	}
}

// m是map类型,用于存储KV数据,filter成员对应每个watch调用时定义的过滤器函数列表.而mu成员为互斥锁,用于在多个goroutine访问或修改时对其他成员提供保护.

// 然后就是Get和Set方法:
func (p *KVStoreService) Get(key string, value,*string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if v, ok := p.m[key]; ok {
		*value = v
		return nil
	}

	return fmt.Errorof("not found")
}

func (p *KVStoreService) Set(kv [2]string, reply *struct{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	key, value := kv[0], kv[1]

	if oldValue := p.m[key]; oldValue != value {
		for _, fn := range p.filter {
			fn(key)
		}
	}

	p.m[key] = value
	return nil
}

// set方法中输入参数是key和value组成的数组,用一个匿名的空结构体表示忽略了输出参数.当修改某个key对应的值时会调用每一个过滤器函数
// 而过滤器列表在watch方法中提供:

func (p *KVStoreService) Watch(timeoutSecond int, keyChanged *string) error {
	id := fmt.Sprintf("watch-%s-%03d", time.Now(), rand.Int())
	ch := make(chan string, 10) 

	p.mu.Lock()
	p.filter[id] = func(key string) {ch <- key}
	p.mu.Unlock()

	select {
	case <- time.After(time.Duration(timeoutSecond) *time.Second):
		return fmt.Errorof("timeout")
	case key := <- ch:
		*keyChanged = key
		return nil
	}

	return nil
}

// watch方法的输入参数是超时的秒数,当有key变化时将key作为返回值返回,如果超时时间后依然没有key被修改,则返回超时的错误.
// watch的实现中,用唯一的id表示每个watch调用,然后根据id将自身对应的过滤器函数注册到p.filter列表
// 从客户端使用watch方法:

func doClientWork(client *rpc.Client) {
	go func() {
		var keyChanged string
		err := client.Call("KVStoreService.Watch",30,&keyChanged)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("watch:", keyChanged)
	}()

	err := client.Call(
		"KVStoreService.Set",[2]string{"abc", "abc-value"},
		new(struct{}),
	)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 3)
}

// 启动一个独立的Goroutine监控key的变化,同步的watch调用会阻塞,直到有key发生变化或者超时.然后通过set方法修改KV值时,服务器会将变化的key通过watch方法返回.
