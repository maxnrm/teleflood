package main

import (
	"github.com/maxnrm/teleflood/internal/flooder"
	"github.com/maxnrm/teleflood/internal/provider"
)

func main() {
	p := provider.New()
	f := flooder.New(p)

	f.Start()
}
