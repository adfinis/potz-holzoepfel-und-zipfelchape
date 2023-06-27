package tracer

/*
 * opentracing adapter for MongoDB
 *
 * mostly pinched from https://jira.mongodb.org/browse/GODRIVER-739
 */

import (
	"context"
	"sync"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"go.mongodb.org/mongo-driver/event"
)

type tracer struct {
	spans sync.Map
}

func NewTracer() *tracer {
	return &tracer{}
}

const prefix = "mongodb."

func (t *tracer) HandleStartedEvent(ctx context.Context, evt *event.CommandStartedEvent) {
	if evt == nil {
		return
	}
	span, _ := opentracing.StartSpanFromContext(ctx, prefix+evt.CommandName)
	ext.DBType.Set(span, "mongo")
	ext.DBInstance.Set(span, evt.DatabaseName)
	ext.DBStatement.Set(span, string(evt.Command))
	span.SetTag("db.host", evt.ConnectionID)
	ext.SpanKind.Set(span, ext.SpanKindRPCClientEnum)
	ext.Component.Set(span, "golang-mongo")
	t.spans.Store(evt.RequestID, span)
}

func (t *tracer) HandleSucceededEvent(ctx context.Context, evt *event.CommandSucceededEvent) {
	if evt == nil {
		return
	}
	if rawSpan, ok := t.spans.Load(evt.RequestID); ok {
		defer t.spans.Delete(evt.RequestID)
		if span, ok := rawSpan.(opentracing.Span); ok {
			defer span.Finish()
			span.SetTag(prefix+"reply", string(evt.Reply))
			span.SetTag(prefix+"duration", evt.Duration)
		}
	}
}

func (t *tracer) HandleFailedEvent(ctx context.Context, evt *event.CommandFailedEvent) {
	if evt == nil {
		return
	}
	if rawSpan, ok := t.spans.Load(evt.RequestID); ok {
		defer t.spans.Delete(evt.RequestID)
		if span, ok := rawSpan.(opentracing.Span); ok {
			defer span.Finish()
			ext.Error.Set(span, true)
			span.SetTag(prefix+"duration", evt.Duration)
			span.LogFields(log.String("failure", evt.Failure))
		}
	}
}
