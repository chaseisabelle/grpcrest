package main

import (

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
"github.com/golang/protobuf/protoc-gen-go/generator"
)

	func main() {
		//buf, _ := os.ReadFile("/Users/chase/Repositories/grpcrest/api/service.yml")
		//doc, _ := libopenapi.NewDocument(buf)
		//v3m, _ := doc.BuildV3Model()
		fdp := &descriptor.FileDescriptorProto{}

		//for _, pth := range v3m.Model.Paths.PathItems {
		//	if pth.Get != nil {
		//		tmp := "foo"
		//		fdp.Service = append(fdp.Service, &descriptorpb.ServiceDescriptorProto{
		//			Name: &tmp,
		//		})
		//		println(fdp.String())
		//	}
		//}

		fdp.Name = proto.String("name")
		fdp.Service = append(fdp.Service, &descriptorpb.ServiceDescriptorProto{
			Name: proto.String("service-name"),
		})

		b, _ := proto.Marshal(fdp)
		println(generator.New().)
	}
