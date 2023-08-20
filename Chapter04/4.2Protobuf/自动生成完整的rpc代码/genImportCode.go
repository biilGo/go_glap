// 自定义genImportCode方法中生成导入包的代码：

func (p *netrpcPlugin) genImportCode(file *generator.FileDescriptor) {
	p.P(`import "net/rpc`)
}

// 在自定义的genServiceCode方法中为每个服务生成相关的代码，分析可以发现每个服务最重要的是服务的名字，然后每个服务有一组方法。而对于服务定义的方法，最重要的是方法的名字，还有输入参数和输出参数类型的名字。
// 定义一个ServiceSpec类型，用于描述服务的元信息：

type ServiceSpec struct {
	ServiceName string
	MethodList  []ServcieMethodSpec
}

type ServiceMethodSpec struct {
	MethodName     string
	InputTypeName  string
	OutputTypeName string
}

// 新建一个buildServiceSpec方法用来解析每个服务的ServiceSpec元信息：

func (p *netrpcPlugin) buildServceiSpec(
	svc *descriptor.ServiceDescriptorProto,
) *ServceiSpec {
	sepc := &ServiceSpec{
		ServiceName: generator.CamelCase(svc.GetName()),
	}

	for _, m := range svc.Method {
		spec.MethodList = append(spec.MethodList, ServiceMethodSpec{
			MethodName:     generator.CamelCase(m.GetName()),
			InputTypeName:  p.TypeName(p.ObjectNamed(m.GetInputType)),
			OutputTypeName: p.TypeName(p.ObjectNamed(m.GetOutputType())),
		})
	}

	return spec
}

/*
输入参数*descriptor.ServcieDescriptorProto类型,完整描述了一个服务的所有信息.然后通过svc.GetName()就可以获取Protobuf文件中定义的服务的名字.
Protobuf文件中的名字转为Go语言的名字后,需要通过generator.CamelCase函数进行一次转换.
类似的,在for循环中我们通过m.GetName()获取方法的名字,然后再转为Go语言中对应的名字.比较复杂的是对输入和输出参数名字的解析:
首先需要通过m.GetInputType()获取输入参数的类型,然后通过p.ObjectName获取类型对应的类对象的信息,最后获取类对象的名字.
*/

// 然后我们就可以基于buildServceiSpec方法构造的服务的元信息生成服务的代码:

func (p *netrpcPlugin) genServiceCode(svc *descriptor.ServiceDescriptorProto) {
	sprc := p.buildServiceSpec(svc)

	var buf bytes.Buffer
	t := template.Must(template.New("").Parse(templService))
	err := t.Execute(&buf, spec)
	if err != nil {
		log.Fatal(err)
	}

	p.P(buf.String())
}

// 为了便于维护,基于Go语言的模板来生成服务的代码,其中tmplService是服务的模板

type HelloServiceInterface interface {
	Hello(in String, out *String) error
}

func RegisterHelloService(srv *rpc.Server, x HelloService) error {
	if err := srv.RegisterName("HelloService", x); err != nil {
		return err
	}

	return nil
}

type HelloServiceClient struct {
	*rpc.Client
}

var _ HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService(network, address string) (*HelloServcieClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HelloServiceClient{Client: c}, nil
}

func (p *HelloServiceClient) Hello(in String, out *String) error {
	return p.Client.Call("HelloService.Hello", in, out)
}

// 其中HelloService是服务名字,同时还有一系列的方法相关的名字.

