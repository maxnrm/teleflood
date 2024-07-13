package flooder

import (
	"context"
	"fmt"

	"github.com/maxnrm/teleflood/internal/provider"
	"github.com/maxnrm/teleflood/internal/sender"
)

type Flooder struct {
	smap map[string]*sender.Sender
	p    provider.Provider
}

func New(s sender.Sender, p provider.Provider) *Flooder {

	return &Flooder{
		p: p,
	}
}

// TODO: implement priority base sending
// priority 1 users who have not been sent messages in N seconds
// priority 2 users who already have been sent messages in N seconds
// priority 3 broadcasted messages
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

		err = s.Send(&msg)
		if err != nil {
			return err
		}
	}
}
