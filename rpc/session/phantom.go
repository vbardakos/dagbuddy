package session

import (
	"context"
)

type phantomID struct{}

type phantomManager[T any] struct {
	newFn func() *T
}

func NewGenPhantomManager[T any](f func() *T) *phantomManager[T] {
	return &phantomManager[T]{newFn: f}
}

func (sm *phantomManager[T]) Attach(ctx context.Context) context.Context {
	return context.WithValue(ctx, phantomID{}, sm.newFn())
}

func GetPhantomFromContext[T any](ctx context.Context) *T {
	return ctx.Value(phantomID{}).(*T)
}
