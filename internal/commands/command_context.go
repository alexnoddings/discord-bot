package commands

import "github.com/bwmarrin/discordgo"

type CommandContext struct {
	Session discordgo.Session
	Event   discordgo.MessageCreate
	Args    []string
}

func (ctx *CommandContext) Reply(message string) (*discordgo.Message, error) {
	msg, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, message)
	return msg, err
}
