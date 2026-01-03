package server

import "context"

type doneNowID struct{}
type doneAfterID struct{}

// immediately closes rpc connection
func Done(ctx context.Context) {
	if fn, ok := ctx.Value(doneNowID{}).(func()); ok {
		fn()
	}
}

// closes rpc connection after dispatch
func DoneAfterDispatch(ctx context.Context) {
	if fn, ok := ctx.Value(doneAfterID{}).(func()); ok {
		fn()
	}
}

func attachDoneNowContextValue(ctx context.Context, ch chan struct{}) context.Context {
	return context.WithValue(ctx, doneNowID{}, buildDoneStatus(ch))
}

func attachDoneAfterContextValue(ctx context.Context, ch chan struct{}) context.Context {
	return context.WithValue(ctx, doneAfterID{}, buildDoneStatus(ch))
}

func buildDoneStatus(ch chan struct{}) func() {
	return func() {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
