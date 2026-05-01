package main

import (
	"context"
	"log"
	"sync"
	"time"
)

// JOBを処理するworkerを作成
func worker(ctx context.Context, id int, jobch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobch:
			if !ok {
				log.Printf("worker %d: job channel closed\n", id)
				return
			}
			log.Printf("worker %d: processing job: %s\n", id, job)
			select {
			// 処理したと仮定して2秒待機する
			case <-time.After(2 * time.Second):
				log.Printf("worker %d: finished\n", id)
			// 処理中にキャン背ルされた場合、ctx.Done()が閉じるので、workerを終了する
			case <-ctx.Done():
				log.Printf("worker %d: canceled while processing\n", id)
				return
			}
		case <-ctx.Done():
			log.Printf("worker %d: stopped\n", id)
			return
		}
	}
}

// JOBを送信するプロデューサーを作成
func enqueue(ctx context.Context, jobch chan<- string, jobs string) error {
	select {
	case jobch <- jobs:
		log.Printf("sent job: %s\n", jobs)
		return nil
	case <-ctx.Done():
		log.Println("send cancelled")
		return ctx.Err()
	}
}

func main() {
	// デフォルトのコンテキストを作成
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 文字列を受け取るチャネルを作成
	jobch := make(chan string, 5)
	// waitGroupを作成して、goroutineの終了を待つ

	var wg sync.WaitGroup
	workerCount := 3
	// workerCountの数だけworkerを起動する
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(ctx, i, jobch, &wg)
	}

	// 仮のJOBを作成する
	jobs := []string{
		"post created: 1",
		"post created: 2",
		"post created: 3",
		"post created: 4",
		"post created: 5",
	}

	// jobchにデータを送信する
	for _, job := range jobs {
		err := enqueue(ctx, jobch, job)
		if err != nil {
			log.Printf("error sending job: %v\n", err)
			break
		}
	}

	close(jobch)
	log.Println("producer closed job channel")

	// waitGroupでworkerの終了を待つ
	wg.Wait()
	log.Println("all workers stopped")
}
