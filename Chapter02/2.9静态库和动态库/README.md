# go_glap
Go-language-advanced-programming

# 静态库和动态库
CGO在使用C/C++资源的时候一般有三种形式:直接使用源码;链接静态库;链接动态库.直接使用源码就是在`import "C"`之前的注释部分包含C代码,或者在当前包中包含C/C++源文件.链接静态库和动态库的方式比较类似,都是通过在LDFLAGS选项指定要链接的库方式链接.

## 使用C静态库
如果CGO中引入C/C++资源有代码而且代码规模也比较小,直接使用源码是最理想的方式,但很多时候我们并没有源代码,或者从C/C++源代码构建的过程异常复杂,这种时候使用C静态库也是一个不错的选择.静态库因为是静态链接,最终的目标程序并不会产生额外的运行时依赖,也不会出现动态库特有的跨运行时资源管理的错误.不过静态库对链接阶段会有一定要求:静态库一般包含了全部的代码,里面会有大量的符号,如果不同静态库之间出现了符号冲突则会导致链接的失败.

先用纯C语言构造一个简单的静态库.要构造的静态库名叫number,库中只有一个number_add_mod函数,用于表示数论中的模加法运算.number库的文件都在number目录下.

`number/number.h`头文件只有一个纯C语言风格的函数声明:
`int number_add_mod(int a, int b, int mod);`

`number/number.c`对应函数的实现:
```c
#include "number.h"
int number_add_mod(int a, int b, int mod) {
    return (a+b)%mod;
}
```

因为CGO使用的是GCC命令来编译和链接C和Go桥接的代码.因此静态库也必须是GCC兼容的格式.

通过以下命令可以生成一个叫libnumber.a的静态库:
```c
$ cd ./number
$ gcc -c -o number.o number.c
$ ar rcs libnumber.a number.o
```

生成libnumber.a静态库之后,就可以在CGO中使用该资源了.

创建main.go文件:
```go
package main

//#cgo CFLAGS: -I./number
//#cgo LDFLAGS: -L${SRCDIR}/number -lnumber
//
//#include "number.h"

import "C"
import "fmt"
func main() {
    fmt.Println(C.number_add_mod(10, 5, 12))
}
```

其中有两个#cgo命令,分别是编译和链接参数.CFLAGS通过`-I./number`将number库对应头文件所在的目录加入头文件检索路径.LFGLAGS通过`-L${SRCDIR}/number`将编译后number静态库所在目录加为链接库检索路径,`-Lnumber`表示链接libnumber.a静态库.需要注意的是,在链接部分的检索路径不能使用相对路径,必须通过cgo特有的`${SRCDIR}`变量将源文件对应的当前目录路径展开为绝对路径.

因为我们有number库的全部代码,所以我们用go generate工具来生成静态库,或者是通过Makefile来构建静态库.因此发布CGO源码包时,我们并不需要提前构建C静态库.

因为多了一个静态库的构建步骤,这种使用了自定义静态库并已经包含了静态库全部代码的Go包无法直接用go get安装.不过我们依然可以通过go get下载,然后用go generate触发静态库构建,最后才是go instann来完成安装.

为了支持go get命令直接下载安装,C语言的`#include`语法可以将number库的源文件链接到当前的包.

创建`z_link_number_c.c`文件如下:
`#include "./number/number.c"`

然后在执行go get或go build之类的命令的时候,CGO就是自动构建number库对应的代码.这种技术是在不改变静态库源代码组织结构的前提下,将静态库转化为了源代码方式引用.

如果使用的第三方的静态库,我们需要先下载安装静态库到合适的位置,然后在#cgo命令中通过CFLAGS和LDFLAGS来指定头文件和库的位置.对于不同的操作系统甚至同一种操作系统的不同版本来说,这些库的安装路径可能都是不同的,那么如何在代码中指定这些可能的变化的参数?

Linux环境,有一个pkg-config命令可以查询要使用某个静态库或动态库时的编译和链接参数.我们可以在#cgo命令中直接使用pkg-config命令来生成编译和链接参数.而且还可以通过pkg_config环境变量定制pkg-config命令.因为不同的操作系统对pkg-config命令的支持不尽相同,通过该方式很难兼容不同的操作系统下的构建参数.不过对于linux等特定的系统,pkg-config命令确实可以简化构建参数的管理.

## 使用C动态库
动态库出现的初衷是对于相同的库,多个进程可以共享同一个,以节省内存和磁盘资源.但是在磁盘和内存已经白菜价的今天,这两个作用已经显得微不足道,那么除此之外动态库还有哪些存在的价值?从库开发角度来说,动态库可以隔离不同动态库之间的关系,减少链接时出现符号冲突的风险.而且对于windows等平台,动态库是跨越VC和GCC不同编译器平台的唯一的可行方式.

对于CGO来说,使用动态库和静态库是一样的,因为动态库也必须是有一个小的静态导出库用于链接动态库.

对于macOS和Linux系统下的gcc环境,我们可以用以下命令创建number库的动态库:
```
$ cd number
$ gcc -shared -o libnumber.so number.c
```

因为动态库和静态库的基础名称都是libnumber,只是后缀不同而已.因此Go语言部分的代码和静态库版本完全一样:
```go
package main

//#cgo CFLAGS: -I./number
//#cgo LDFLAGS: -L${SRCDIR}/number -lnumber
//
//#include "number.h"
import "C"
import "fmt"

func main() {
    fmt.Println(C.number_add_mod(10, 5, 12))
}
```

编译时GCC会自动找到libnumber.a或libnumber.so进行链接.

windows平台,还可以用VC工具来生成动态库.需要先为number.dll创建一个def文件,用于控制要导出到动态库的符号.

number.def文件的内容如下:
```
LIBRARY number.dll
EXPORTS
number_add_mod
```

其中第一行的LIBRARY指明动态库的文件名,然后的EXPORTS语句之后是要导出的符号列表.

现在我们可以用以下命令来创建动态库.
```
$ cl /c number.c
$ link /DLL /OUT:number.dll number.obj number.def
```

这时候会为dll同时生成一个number.lib的导出库.但是CGO中我们无法使用lib格式的链接库.

要生成`.a`格式的导出库需要通过mingw工具箱中的dlltool命令完成:
`$ dlltool -dllname number.dll --def number.def --output-lib libnumber.a`

生成了libnumber.a文件之后,就可以通过`-lnumber`链接参数进行链接了.

需要注意的是,在运行时需要将动态库放到系统能够找到的位置.

## 导出C静态库
CGO不仅可以使用C静态库,也可以将Go实现的函数导出为C静态库.

创number.go:
```go
package main

import "C"

func main() {}

//export number_add_mod
func number_add_mod(a, b, mod C.int) C.int {
    return (a + b) % mod
}
```

根据CGO文档的要求,我们需要在main包中导出C函数,对于C静态库构建方式来说,会忽略main包中的main函数,只是简单导出C函数,采用以下命令构建:
`$ go build -buildmode=c-archive -o number.a`

在生成number.a静态库的同时,cgo还会生成一个number.h文件.

number.h文件的内容如下:
```
#ifdef __cplusplus
extern "C" {
#endif
extern int number_add_mod(int p0, int p1, int p2);
#ifdef __cplusplus
}
#endif
```

其中`extern "C"`部分的语法是为了同时适配C和C++两种语言.核心内容是声明了要导出的number_add_mod函数.

然后我们创建了一个`_test_main.c`的C文件用于测试生成的C静态库:
```c
#include "number.h"
#include <stdio.h>
int main() {
    int a = 10;
    int b = 5;
    int c = 12;
    int x = number_add_mod(a, b, c);
    printf("(%d+%d)%%%d = %d\n", a, b, c, x);
    return 0;
}
```

通过以下命令编译并运行:
```
$ gcc -o a.out _test_main.c number.a
$ ./a.out
```

## 导出C动态库
CGO导出动态库的过程和静态库类似,只是将构建模式改为`c-shared`,输出文件名改为`number.so`而已:
`$ go build -buildmode=c-shared -o number.so`

`_test_main.c`文件内容不变,然后用以下命令编译并运行:
```
$ gcc -o a.out _test_main.c number.so
$ ./a.out
```

## 导出非main包的函数
通过`go help buildmode`命令可以查看C静态库和C动态库的构建说明:
```
-buildmode=c-archive
    Build the listed main package, plus all packages it imports,
    into a C archive file. The only callable symbols will be those
    functions exported using a cgo //export comment. Requires
    exactly one main package to be listed.
-buildmode=c-shared
    Build the listed main package, plus all packages it imports,
    into a C shared library. The only callable symbols will
    be those functions exported using a cgo //export comment.
    Requires exactly one main package to be listed.
```

文档说明导出的C函数必须是在main包导出,然后才能在生成的头文件包含声明的语句,但是很多时候我们可能更希望将不同类型的导出函数组织到不同的Go包中,然后统一导出为一个静态库或动态库.

要实现从是从非main包导出C函数,或者是多个包导出C函数,我们需要自己提供导出C函数对应的头文件.

假设我们先创建一个number子包,用于提供摸加法函数:
```go
package number
import "C"
//export number_add_mod
func number_add_mod(a, b, mod C.int) C.int {
    return (a + b) % mod
}
```

然后是当前的main包:
```go
package main
import "C"
import (
    "fmt"
    _ "./number"
)
func main() {
    println("Done")
}
//export goPrintln
func goPrintln(s *C.char) {
    fmt.Println("goPrintln:", C.GoString(s))
}
```

其中我们导入了number子包,在number子包中有导出的C函数number_add_mod,同时我们在main包也导出了goPrintln函数

通过以下命令创建C静态库:
`$ go build -buildmode=c-archive -o main.a`

这时候在生成main.a静态库的同时,也会生成一个main.h头文件.但是main.h头文件中只有main包导出的goPrintln函数的声明,并没有number子包导出函数的声明.其实number_add_mod函数在生成的C静态库中是存在的,我们可以直接使用.

创建`_test_main.c`测试文件如下:
```c
#include <stdio.h>
void goPrintln(char*);
int number_add_mod(int a, int b, int mod);
int main() {
    int a = 10;
    int b = 5;
    int c = 12;
    int x = number_add_mod(a, b, c);
    printf("(%d+%d)%%%d = %d\n", a, b, c, x);
    goPrintln("done");
    return 0;
}
```

我们并没有包含CGO自动生成的main.h头文件,而是通过手动方式声明了goPrintln和number_add_mod两个导出函数.