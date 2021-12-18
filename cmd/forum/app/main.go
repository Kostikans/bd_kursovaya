package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kostikan/bd_kursovaya/internal/app/api/forum"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"

	f "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
)

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pkg/swagger/api/forum/forum.swagger.json")
}

func startGRPC() {
	lis, err := net.Listen("tcp", "localhost:8079")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	forum := forum.NewForum(forum.Opts{})
	f.RegisterForumServer(grpcServer, forum)
	log.Println("gRPC server ready...")
	fmt.Println(grpcServer.Serve(lis))
}

func startHTTP() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Connect to the GRPC server
	conn, err := grpc.Dial("localhost:8079", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	// Register grpc-gateway
	rmux := runtime.NewServeMux()
	client := f.NewForumClient(conn)
	err = f.RegisterForumHandlerClient(ctx, rmux, client)
	if err != nil {
		log.Fatal(err)
	}

	// Serve the swagger,
	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	mux.HandleFunc("/docs/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/forum.swagger.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))

	mux.Handle("/forum.swagger.json", http.FileServer(http.Dir("pkg/swagger/api/forum")))

	log.Println("REST server ready...")
	log.Println("Serving Swagger at: http://localhost:8080/docs/")
	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	go startGRPC()
	time.Sleep(1 * time.Second)
	go startHTTP()

	// Block forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()

}
