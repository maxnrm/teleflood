package main

import (
	"fmt"

	"github.com/maxnrm/teleflood/internal/flooder"
	"github.com/maxnrm/teleflood/internal/provider"
)

func main() {
	p := provider.New()
	f := flooder.New(p)

	fmt.Println("started")
	f.Start()
}
