package flooder

import (
	"context"
	"fmt"

	"github.com/maxnrm/teleflood/config"
	"github.com/maxnrm/teleflood/internal/provider"
	"github.com/maxnrm/teleflood/internal/sender"
	"go.uber.org/ratelimit"
)

type Flooder struct {
	grl ratelimit.Limiter
	// TODO: delete senders unused for some time
	smap map[string]*sender.Sender
	p    provider.Provider
}

func New(p provider.Provider) *Flooder {

	grl := ratelimit.New(config.GLOBAL_RATE_LIMIT_GLOBAL)

	return &Flooder{
		p:   p,
		grl: grl,
	}
}

// TODO: implement priority base sending
// priority 1 users who have not been sent messages in N seconds
// priority 2 users who already have been sent messages in N seconds
// priority 3 broadcasted messages

// TODO 2: more flexible ratelimit to be able to return tokens if not used,
// OR some other support of balancing bulk message to multiple users and tokens
func (f *Flooder) Start() error {
	for {
		ctx := context.Background()
		wrappedMsg, err := f.p.Next(ctx)
		if err != nil {
			return err
		}

		token := wrappedMsg.BotToken
		msg := wrappedMsg.FloodMessage

		s, ok := f.smap[token]
		if !ok {
			s, err = sender.New(token)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		go s.Send(f.grl, &msg)
	}
}
