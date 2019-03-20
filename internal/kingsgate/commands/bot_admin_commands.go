package commands

import (
	"kingsgate/internal/commands"
	"kingsgate/internal/commands/predicates"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var (
	botCommand = commands.Command{
		Name:               "Bot",
		Aliases:            []string{"Admin"},
		RequiredPredicates: []func(ctx *commands.CommandContext) error{predicates.IsBotOwnerPredicate},
		OnInvoked:          onBotInvoked,
	}

	botShutdownCommand = commands.Command{
		Name:      "Shutdown",
		Aliases:   []string{"End", "Kill"},
		OnInvoked: onBotShutdownInvoked,
	}

	botStatusCommand = commands.Command{
		Name:      "Status",
		OnInvoked: onBotStatusInvoked,
	}
)

func onBotInvoked(ctx *commands.CommandContext) error {
	return errors.New("no sub-command invoked")
}

func onBotShutdownInvoked(ctx *commands.CommandContext) error {
	os.Exit(1)
	return nil
}

func onBotStatusInvoked(ctx *commands.CommandContext) error {
	newStatus := strings.Join(ctx.Args, " ")
	return ctx.Session.UpdateStatus(0, newStatus)
}

func createBotAdminCommand() *commands.Command {
	botCommand.AddChild(&botShutdownCommand)
	botCommand.AddChild(&botStatusCommand)

	return &botCommand
}
