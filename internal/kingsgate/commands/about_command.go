package commands

import (
	"fmt"
	"kingsgate/internal/commands"
	"os"
	"runtime"
)

var (
	helpCommand = commands.Command{
		Name:      "Help",
		OnInvoked: onHelpInvoked,
	}
	helpAboutCommand = commands.Command{
		Name:      "About",
		Aliases:   []string{"Info"},
		OnInvoked: onHelpAboutInvoked,
	}
)

func onHelpInvoked(ctx *commands.CommandContext) error {
	ctx.Reply("```" + commandHandler.GenerateTreeView() + "```")
	return nil
}

func onHelpAboutInvoked(ctx *commands.CommandContext) error {
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)

	host, _ := os.Hostname()
	goos := runtime.GOOS
	arch := runtime.GOARCH

	alloc := memory.Alloc / 1048576
	totalAlloc := memory.TotalAlloc / 1048576
	sys := memory.Sys / 1048576

	ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, fmt.Sprintf(
		"```asciidoc\n"+
			"[ Discord Bot running Kingsgate ]\n"+
			"\n"+
			"= Runtime =\n"+
			"  System       :: %s\n"+
			"  OS Target    :: %s\n"+
			"  Architecture :: %s\n"+
			"= Memory =\n"+
			"  Sys        :: %dMB\n"+
			"  Alloc      :: %dMB\n"+
			"  TotalAlloc :: %dMB\n"+
			"```",
		host, goos, arch,
		sys, alloc, totalAlloc))
	return nil
}

func createAboutCommand() *commands.Command {
	helpCommand.AddChild(&helpAboutCommand)

	return &helpCommand
}
