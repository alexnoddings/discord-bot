package commands

import (
	"kingsgate/internal/commands"

	"github.com/bwmarrin/discordgo"
)

const (
	prefix = "."
)

var (
	commandHandler commands.CommandHandler
)

// Register registers a command handler and root commands to a session
func Register(session *discordgo.Session) {
	commandHandler = commands.NewCommandHandler(prefix)

	commandHandler.AddRootCommand(createAboutCommand())
	commandHandler.AddRootCommand(createAudioCommand())
	commandHandler.AddRootCommand(createBotAdminCommand())

	session.AddHandler(commandHandler.OnMessageCreated)
}
