package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/dig"
	"google.golang.org/grpc"

	f "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
)

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pkg/swagger/api/forum/forum.swagger.json")
}

// MainOpts - main options
type MainOpts struct {
	Components map[string]interface{}
}

// Main - main
func Main(params ...func(*MainOpts)) {
	ctx := context.Background()
	opts := &MainOpts{
		Components: map[string]interface{}{
			"forum service": forumServiceProvider,
			"facade":        facadeProvider,

			"database":     databaseProvider,
			"account repo": accountRepoProvider,
			"comment repo": commentRepoProvider,
			"post repo":    postRepoProvider,
			"tag repo ":    tagRepoProvider,
		},
	}
	for _, param := range params {
		param(opts)
	}

	c := dig.New()

	err := provideAndInvoke(ctx, c, opts.Components, func(ac appContext) {
		lis, err := net.Listen("tcp", "localhost:8079")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		f.RegisterForumServer(grpcServer, ac.Forum)
		log.Println("gRPC server ready...")
		go grpcServer.Serve(lis)

		time.Sleep(1 * time.Second)
		go startHTTP()

		var wg sync.WaitGroup
		wg.Add(1)
		wg.Wait()
	})
	if err != nil {
		log.Fatal(err)
	}
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
