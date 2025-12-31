package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/NebuxCloud/botbuster/internal/api"
	"github.com/NebuxCloud/botbuster/internal/crypto"
	"github.com/NebuxCloud/botbuster/internal/services"
	"github.com/samber/do/v2"
)

func main() {
	ctx := context.Background()

	// Create the dependency injection container
	i := do.New(services.Package)

	defer i.ShutdownWithContext(ctx)
	go i.ShutdownOnSignalsWithContext(ctx)

	// Get command argument
	arg := ""

	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	// Run command
	log := do.MustInvoke[*slog.Logger](i)

	var err error

	switch arg {
	case "", "serve":
		api := do.MustInvoke[*api.API](i)
		err = api.Serve(ctx)

		if err != nil {
			log.Error(err.Error())
		}
	case "generate:key":
		key, err := crypto.GenerateHMACKey()

		if err != nil {
			log.Error(err.Error())
		}

		fmt.Printf("%s\n", key)
	}
}
