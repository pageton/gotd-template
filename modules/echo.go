package modules

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type EchoModule struct{}

var echoRegex = regexp.MustCompile(`^/echo(?:\s+(.+))?$`)

func (e *EchoModule) Register(
	ctx context.Context,
	sender *message.Sender,
	dispatcher *tg.UpdateDispatcher,
) {
	dispatcher.OnNewMessage(
		func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage) error {
			m, ok := u.Message.(*tg.Message)
			if !ok || m.Out {
				return nil
			}

			matches := echoRegex.FindStringSubmatch(strings.TrimSpace(m.Message))
			if matches == nil {
				return nil
			}

			log.Println("Echo command received")
			response := "Echo!"
			if len(matches) > 1 && matches[1] != "" {
				response = matches[1]
			}

			_, err := sender.Reply(entities, u).Text(ctx, response)
			return err
		},
	)
}
