package main

import (
	"context"
	"log"
	"sync"
)

func main() {
	// デフォルトのコンテキストを作成
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// waitGroupを作成して、goroutineの終了を待つ
	var wg sync.WaitGroup
	wg.Add(1)

	// 文字列を受け取るチャネルを作成
	jobch := make(chan string, 100)

	// jobchから文字列を受け取る関数を作成してgoルーチンで実行 ※最後の()で実行する
	go func() {
		defer wg.Done()
		for {
			select {
			case job, ok := <-jobch:
				// チャネルがcloseされた場合、ok=falseになるのでworkerを終了する
				// close後はゼロ値が返り続けるため、okチェックしないと無限ループになる
				if !ok {
					log.Printf("job channel closed")
					return
				}
				log.Printf("received job: %s", job)
			case <-ctx.Done():
				log.Printf("worker stopped")
				return
			}
		}
	}()

	// jobchにデータを送信する
	for range 100 {
		select {
		case jobch <- "post created: 123":
			log.Println("job sent")
		case <-ctx.Done():
			log.Println("send cancelled")
			close(jobch)
			return
		}
	}

	close(jobch)
	// waitGroupでworkerの終了を待つ
	wg.Wait()
	log.Println("main stopped")
}
