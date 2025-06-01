package modules

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

func InitModules(ctx context.Context, api *tg.Client, dispatcher tg.UpdateDispatcher) {
	api.UpdatesGetState(ctx)
	sender := message.NewSender(api)

	modules := []Module{
		&EchoModule{},
	}

	for _, m := range modules {
		m.Register(ctx, sender, &dispatcher)
	}
}
