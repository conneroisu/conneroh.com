package cache

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

// Querer reads from the message channel and sends messages to the que channel.
func Querer(
	ctx context.Context,
	queCh MsgChannel,
	msgCh MsgChannel,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-queCh:
			slog.Info("queing", "msg", msg)
			msgCh <- msg
		}
	}
}

// CtxSend sends a message to the channel with a context.
func CtxSend[T any](
	ctx context.Context,
	ch chan T,
	msg T,
) error {
	select {
	case <-time.After(time.Second * 5):
		err := errors.New("timeout sending message")
		slog.Error(err.Error(), "msg", msg)
		return err
	case <-ctx.Done():
		return ctx.Err()
	case ch <- msg:
		return nil
	}
}
