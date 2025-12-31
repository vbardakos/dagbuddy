package server

import (
	"context"

	p "github.com/vbardakos/dagbuddy/rpc/protocol"
)

type dispatcher struct {
	selectFunc Selector
	out        chan<- p.RPCMessage
	onError    func(error)
}

func (d *dispatcher) dispatch(ctx context.Context, msg p.RPCMessage) {
	handler := d.selectFunc(msg)

	if handler == nil {
		return
	}

	go func() {
		r, err := handler.Handle(ctx, msg)

		if err != nil {
			d.onError(err)
			return
		}

		if r.Type() == p.NotificationKind {
			return
		}

		select {
		case <-ctx.Done():
			return

		case d.out <- r:
		}
	}()
}
