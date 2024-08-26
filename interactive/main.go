package interactive

import (
	"google.golang.org/grpc"
	"log"
	"net"
	intrv1 "we_book/api/proto/gen/intr"
	grpc1 "we_book/interactive/grpc"
)

func main() {
	server := grpc.NewServer()
	intrServer := &grpc1.InteractiveServiceServer{}
	intrv1.RegisterInteractiveServiceServer(server, intrServer)
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	err = server.Serve(l)
	log.Println(err)
}
