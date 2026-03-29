package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

// errgroupを使った並列処理サンプル
// 1つのgoroutineがエラーを返すと、contextがcancelされ他のgoroutineも終了する
func main() {
	parent := context.Background()
	g, ctx := errgroup.WithContext(parent)

	g.Go(func() error {
		time.Sleep(2 * time.Second)
		fmt.Println("worker1 done")
		return nil
	})

	g.Go(func() error {
		time.Sleep(1 * time.Second)
		return fmt.Errorf("worker2 failed")
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			fmt.Println("worker3 canceled")
			return ctx.Err()
		case <-time.After(5 * time.Second):
			fmt.Println("worker3 done")
			return nil
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
