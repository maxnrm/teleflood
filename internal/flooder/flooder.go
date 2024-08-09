package flooder

import (
	"context"
	"fmt"
	"time"

	"github.com/maxnrm/teleflood/config"
	"github.com/maxnrm/teleflood/internal/sender"
	"github.com/maxnrm/teleflood/pkg/message"
	"golang.org/x/time/rate"
)

type provider interface {
	Next(context.Context) (*message.WrappedMessage, error)
}

type Flooder struct {
	grl *rate.Limiter
	// TODO: delete senders unused for some time
	smap map[string]*sender.Sender
	p    provider
}

func New(p provider) *Flooder {

	limit := rate.Every(time.Second / time.Duration(config.GLOBAL_RATE_LIMIT_GLOBAL))

	return &Flooder{
		p:    p,
		smap: make(map[string]*sender.Sender),
		grl:  rate.NewLimiter(limit, config.GLOBAL_RATE_LIMIT_GLOBAL),
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
			fmt.Println(err)
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
			f.smap[token] = s
		}

		go s.Send(ctx, f.grl, &msg)
	}
}
