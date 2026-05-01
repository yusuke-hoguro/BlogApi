package handler

import (
	"context"
	"log"

	"github.com/yusuke-hoguro/BlogApi/internal/workerpool"
)

// 監視イベントをワーカープールに追加するためのユーティリティ関数
func enqueueAuditEvent(ctx context.Context, auditPool *workerpool.AuditWorkerPool, event workerpool.AuditEvent) {
	if auditPool == nil {
		return
	}
	if err := auditPool.Enqueue(ctx, event); err != nil {
		log.Printf("Failed to enqueue audit event: %v", err)
	}
}
