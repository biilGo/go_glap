# go_glap
Go-language-advanced-programming

# 快速入门

## 最简CGO程序
真实的CGO程序一般都比较复杂.不过我们可以由浅入深,一个最简单的CGO程序该是什么样的?要构造一个最简CGO程序,首先要忽视一些复杂的CGO特性,同时要展示CGO程序和纯go程序的差别来.下面是我们构建的最简CGO程序:
```go
package main
import "C"
func main() {
    println("hello cgo")
}
```

代码通过`import C`语句启用CGO特性,主函数只是通过go内置的println函数输出字符串,其中并没有任何和CGO相关的代码.虽然没有调用CGO的相关函数,但是go build命令会在编译和链接阶段启动gcc编译器,这已经是一个完整的CGO程序.

## 基于C标准库函数输出字符串
来看看更简单的版本:
```go
package main
//#include <stdio.h>
import "C"
func main() {
    C.puts(C.CString("Hello, World\n"))
}
```

我们不仅仅通过`import C`语句启用CGO特性,同时包含C语言的`<stdio.h>`头文件.然后通过CGO包的`C.CString`函数将go语言字符串转为C语言字符串,最后调用CGO包的`C.puts`函数向标准输出窗口打印转换后的C字符串.

相比hello world的革命:我们没有在程序退出前释放`C.CString`创建的C语言字符串;还有我们改用`puts`函数直接向标准输出打印,之前是采用`fputs`向标准输出打印.

没有释放使用`C.CString`创建的C语言字符串会导致内存泄露.但是对于这个小程序来说,这样是没有问题的,因为程序退出后操作系统会自动回收程序的所有资源.

## 使用自己的C函数
前面我们使用了标准库中已有的函数.现在我们自定义一个叫`SayHello`的C函数来实现打印,然后从go语言环境中调用这个`SayHello`函数:
```go
package main
/*
#include <stdio.h>
static void SayHello(const char* s) {
    puts(s);
}
*/
import "C"
func main() {
    C.SayHello(C.CString("Hello, World\n"))
}
```

除了`SayHello`函数是我们自己实现的之外,其他的部分和前面的例子基本相似.

我们也可以将`SayHello`函数放到当前目录下的一个C语言源文件中.因为是编写在独立的C文件中,为了允许外部引用,所以需要去掉函数的`static`修饰符
```go
// hello.c
#include <stdio.h>
void SayHello(const char* s) {
    puts(s);
}
```

然后在CGO部分先声明`SayHello`函数,其他部分不变:
```go
package main
//void SayHello(const char* s);
import "C"
func main() {
    C.SayHello(C.CString("Hello, World\n"))
}
```

既然`SayHello`函数已经放到独立的C文件中了,我们自然可以将对应的C文件编译打包为静态或动态文件供使用.如果是以静态库或动态库方式引用`SayHello`函数的话,需要将对应的C源文件移出当前的目录.

## C代码的模块化
在编译过程中,抽象和模块化是将复杂文件简化的通用手段.当代码的语句变多时,我们可以将相似的代码封装到一个个函数中;当程序中的函数变多时,我们将函数拆分到不同的文件或模块中.而模块化编程的核心是面向程序接口编程.

前面的例子中,我们可以抽象一个名为hello的模块,模块的全部接口函数都在hello.h头文件定义:
```c
// hello.h
void SayHello(const char* s);
```

其中只有一个SayHello函数的声明.但是作为hello模块的用户来说,就可以放心地使用SayHello函数,而无需关心函数的具体实现.而作为SayHello函数的实现者来说,函数的实现只要满足头文件中函数的声明的规范即可.下面是SayHello函数的C语言实现,对应hello.c文件:
```c
// hello.c
#include "hello.h"
#include <stdio.h>
void SayHello(const char* s) {
    puts(s);
}
```

在hello.c文件的开头,实现者通过`#include "hello.h"`语句包含SayHello函数的声明,这样可以保证函数的实现满足模块对外公开的接口.

接口文件hello.h是hello模块的实现者和使用者共同的约定,但是该约定并没有要求必须使用C语言来实现SayHello函数.我们也可以用C++语言来重新实现这个C语言函数:
```C
// hello.cpp
#include <iostream>
extern "C" {
    #include "hello.h"
}
void SayHello(const char* s) {
    std::cout << s;
}
```

在C++版本的SayHello函数实现中,通过C++特有的`std::count`输出流输出的字符串,不过为了保证C++语言实现的SayHello函数满足头文件hello.h定义的函数规范,我们需要通过`extern "C"`语句指示该函数的链接符号遵循C语言的规则.

在采用面向C语言的API接口编程之后,我们彻底解放了模块实现者的语言枷锁:实现则可以用任何编程说实现模块,只要种满足公开的API约定即可.我们可以用C语言实现SayHello函数,也可以使用更复杂的C++语言来实现SayHello函数,当然我们也可以用汇编语句甚至go语言来重新实现SayHello函数

## 用go重新实现C函数
其实CGO不仅仅用于go语言中调用C语言函数,还可以用于导出go语言函数给c语言函数调用,前面的例子中,我们已经抽象一个名为hello的模块,模块的全部接口都在hello.h头文件定义:
```c
// hello.h
void SayHello(/*const*/ char* s);
```

现在我们创建一个hello.go文件,用go语言重新实现C语言接口的SayHello函数:
```go
// hello.go
package main
import "C"
import "fmt"
//export SayHello
func SayHello(s *C.char) {
    fmt.Print(C.GoString(s))
}
```

我们通过CGO的`//export SayHello`指令将go语言实现的函数`SayHello`导出为C语言函数,为了适配CGO导出的C语言函数,我们禁止了在函数的声明语句中的const修饰符.需要注意的是,这里其实有两个版本的SayHello函数:一个go语言环境的;另一个是C语言环境的,CGO生成的C语言版本SayHello函数最终会通过桥接代码调用go语言版本的SayHello函数

通过面向C语言接口的编程技术,我们不仅仅解放了函数的实现者,同时也简化的函数的使用者,现在我们可以将SayHello当作一个标准库的函数使用:
```go
package main
//#include <hello.h>
import "C"
func main() {
    C.SayHello(C.CString("Hello, World\n"))
}
```

一切似乎都回到了开始的CGO代码,但是代码内涵更丰富了

## 面向C接口的go编程
在开始的例子中,我们的全部CGO代码都在一个go文件中,通过面向C接口编程的技术将SayHello分别拆分到不同的C文件,而main依然是go文件,再然后,是用go函数重新实现了c语言接口的SayHello函数.但是对于目前的例子来说只有一个函数,要拆分到三个不同的文件确实有些繁琐.

正所谓合久必分,分久必合,我们尝试将例子中的几个文件重新合并到一个go文件中:
```go
package main
//void SayHello(char* s);
import "C"
import (
    "fmt"
)
func main() {
    C.SayHello(C.CString("Hello, World\n"))
}
//export SayHello
func SayHello(s *C.char) {
    fmt.Print(C.GoString(s))
}
```

现在版本的CGO代码中c语言代码的比例已经很少了,但是我们依然可以进一步以go语言的思维来提炼我们的CGO代码.通过分析发现`SayHello`函数中的参数如果可以直接使用go字符串是最直接的.在go1.10中CGO增加了一个`_GoString_`预定义的C语言类型,用来表示go语言字符串,下面是改进后的代码:
```go
// +build go1.10
package main
//void SayHello(_GoString_ s);
import "C"
import (
    "fmt"
)
func main() {
    C.SayHello("Hello, World\n")
}
//export SayHello
func SayHello(s string) {
    fmt.Print(s)
}
```

虽然看起来全部是go语言代码,但是执行的时候先从go语言的`main`函数,到CGO自动生成的C语言版本`SayHello`桥接函数,最后又回到了go语言环境的`SayHello`函数.