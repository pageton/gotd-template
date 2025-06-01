package modules

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type Module interface {
	Register(ctx context.Context, sender *message.Sender, dispatcher *tg.UpdateDispatcher)
}
