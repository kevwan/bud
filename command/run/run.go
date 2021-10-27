package run

import (
	"context"
	"fmt"

	"gitlab.com/mnm/bud/bfs"
)

type Command struct {
	Hot   bool
	Embed bool
}

func (c *Command) Run(ctx context.Context, generators map[string]bfs.Generator) error {
	fmt.Println("running code!", c.Hot, c.Embed)

	// 1. Run the generators
	// 2. go run bud/main.go
	// 3. Wait for changes
	return nil
}