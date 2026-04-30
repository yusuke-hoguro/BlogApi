package workerpool

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

var ErrQueueFull = errors.New("job queue is full")

// 監視イベントの構造体
type AuditEvent struct {
	Action string
	UserID int
	PostID int
}

// 監視ワーカープールの構造体
type AuditWorkerPool struct {
	jobCh       chan AuditEvent
	workerCount int
	wg          sync.WaitGroup
}

// 新規監視ワーカープールの作成
func NewAuditWorkerPool(wokercount int, queueSize int) *AuditWorkerPool {
	return &AuditWorkerPool{
		jobCh:       make(chan AuditEvent, queueSize),
		workerCount: wokercount,
	}
}

// 監視ワーカープールの開始
func (p *AuditWorkerPool) Start(ctx context.Context) {
	for i := 1; i <= p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}
}

// 監視ワーカープールのキューにイベントを追加
func (p *AuditWorkerPool) Enqueue(ctx context.Context, event AuditEvent) error {
	select {
	case p.jobCh <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return ErrQueueFull
	}
}

// 監視ワーカープールの停止
func (p *AuditWorkerPool) Stop() {
	close(p.jobCh)
	p.wg.Wait()
}

// ワーカーの処理ループ
func (p *AuditWorkerPool) worker(ctx context.Context, id int) {
	defer p.wg.Done()

	for {
		select {
		case event, ok := <-p.jobCh:
			if !ok {
				log.Printf("audit worker %d: job channel closed", id)
				return
			}

			// イベントの処理を実施
			if err := processAuditEvent(ctx, id, event); err != nil {
				log.Printf("audit worker %d: failed to process event: %v", id, err)
			}
		case <-ctx.Done():
			log.Printf("audit worker %d: context cancelled", id)
			return
		}
	}
}

// 監視イベントの処理関数
func processAuditEvent(ctx context.Context, workerID int, event AuditEvent) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// ここで実際のイベント処理を行う（例: ログ出力）
		log.Printf("audit worker %d: processing event: action=%s, userID=%d, postID=%d", workerID, event.Action, event.UserID, event.PostID)
		return nil
	}
}

// 監視イベントを文字列に変換する関数
func (e AuditEvent) String() string {
	return fmt.Sprintf("action=%s user_id=%d post_id=%d", e.Action, e.UserID, e.PostID)
}
