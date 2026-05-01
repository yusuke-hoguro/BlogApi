package workerpool

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

var ErrQueueFull = errors.New("job queue is full")
var ErrQueueClosed = errors.New("job queue is closed")

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
	stopOnce    sync.Once
	wg          sync.WaitGroup
	mu          sync.RWMutex
	closed      bool
}

// 新規監視ワーカープールの作成
func NewAuditWorkerPool(wokercount int, queueSize int) *AuditWorkerPool {
	return &AuditWorkerPool{
		jobCh:       make(chan AuditEvent, queueSize),
		workerCount: wokercount,
	}
}

// 監視ワーカープールの開始
func (p *AuditWorkerPool) Start() {
	for i := 1; i <= p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// 監視ワーカープールのキューにイベントを追加
func (p *AuditWorkerPool) Enqueue(ctx context.Context, event AuditEvent) error {
	// 読み取り用のロック取得
	p.mu.RLock()
	defer p.mu.RUnlock()

	// ワーカープールが停止している場合はエラーを返す
	if p.closed {
		return ErrQueueClosed
	}

	// ジョブチャンネルにイベントを送信（非ブロッキング）
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
	p.stopOnce.Do(func() {
		p.mu.Lock()
		if !p.closed {
			p.closed = true
			close(p.jobCh)
		}
		p.mu.Unlock()
		p.wg.Wait()
	})
}

// ワーカーの処理ループ
func (p *AuditWorkerPool) worker(id int) {
	defer p.wg.Done()

	for event := range p.jobCh {
		// イベントの処理を実施
		if err := processAuditEvent(id, event); err != nil {
			log.Printf("audit worker %d: failed to process event: %v", id, err)
		}
	}

	log.Printf("audit worker %d: job channel closed", id)
}

// 監視イベントの処理関数
func processAuditEvent(workerID int, event AuditEvent) error {
	// ここで実際のイベント処理を行う（例: ログ出力）
	log.Printf("audit worker %d: processing event: action=%s, userID=%d, postID=%d", workerID, event.Action, event.UserID, event.PostID)
	return nil
}

// 監視イベントを文字列に変換する関数
func (e AuditEvent) String() string {
	return fmt.Sprintf("action=%s user_id=%d post_id=%d", e.Action, e.UserID, e.PostID)
}
