package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
	"google.golang.org/grpc"
)

const (
	commentReqCount = 10
)

type perf struct {
	workers  []worker
	reqCount uint64
	mtx      *sync.RWMutex
}

type worker struct {
	p *perf
}

func main() {
	var workerCount int
	flag.IntVar(&workerCount, "c", 100, "worker count")
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	conn, err := grpc.Dial("localhost:8079", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	client := desc.NewForumClient(conn)
	perf := &perf{
		mtx:      &sync.RWMutex{},
		reqCount: 0,
	}
	perf.workers = make([]worker, 0, workerCount)
	for i := 0; i < workerCount; i++ {
		perf.workers = append(perf.workers, worker{p: perf})
		go perf.workers[i].loop(ctx, client)
	}

	ticker := time.NewTicker(1000 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:

				return
			case _ = <-ticker.C:
				fmt.Println(fmt.Sprintf("rps: %d", perf.getReqCountAndClear()))
			}
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for {
			select {
			case _ = <-sigc:
				cancel()
				done <- true
				fmt.Println("done")
				os.Exit(0)
			}
		}
	}()
	for {
	}
}

func (w *worker) loop(ctx context.Context, conn desc.ForumClient) {
	for {
		w.performRequests(ctx, conn)
	}
}

func (w *worker) performRequests(ctx context.Context, conn desc.ForumClient) {
	posts, err := conn.GetPosts(ctx, &desc.GetPostListRequest{
		Cursor: 0,
		Limit:  1000,
	})
	if err != nil {
		log.Fatalln(err)
	}
	w.p.incrementRequests()

	postID := posts.GetPosts()[rand.Intn(len(posts.GetPosts()))].GetId()
	for i := 0; i < commentReqCount; i++ {
		_, err := conn.GetComments(ctx, &desc.GetCommentListRequest{
			Cursor: 0,
			Limit:  1000,
			PostId: postID,
		})
		if err != nil {
			log.Fatalln(err)
		}
		w.p.incrementRequests()
	}
}

func (p *perf) incrementRequests() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.reqCount++
}

func (p *perf) getReqCountAndClear() uint64 {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	requests := p.reqCount
	p.reqCount = 0
	return requests
}
