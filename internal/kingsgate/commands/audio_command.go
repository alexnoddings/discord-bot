package commands

import (
	"kingsgate/internal/commands"
	"kingsgate/internal/commands/predicates"
	"kingsgate/internal/kingsgate/audio"
	"kingsgate/internal/kingsgate/config"
	"strings"

	"github.com/pkg/errors"
)

const (
	// Move Members = 0x1000000
	audioPlayUserRequiredPermissions = 0x1000000

	// Manage Roles = 0x10000000, Move Members = 0x1000000
	audioPlayBotRequiredPermissions = 0x10000000 | 0x1000000

	// Manage Roles = 0x10000000, Move Members = 0x1000000, Manage Channels = 0x10
	audioConfigUserRequiredPermissions = 0x10000000 | 0x1000000 | 0x10

	// Entirely bot-sided config
	audioConfigBotRequiredPermissions = 0
)

var (
	audioParentCommand = commands.Command{
		Name:      "Audio",
		OnInvoked: onAudioParentInvoked,
	}

	audioPlayCommand = commands.Command{
		Name: "Play",
		RequiredPredicates: []func(ctx *commands.CommandContext) error{
			predicates.GenerateAuthorHasPermissionPredicate(audioPlayUserRequiredPermissions),
			predicates.GenerateBotHasPermissionPredicate(audioPlayBotRequiredPermissions)},
		OnInvoked: onAudioPlayInvoked,
	}
)

func onAudioParentInvoked(ctx *commands.CommandContext) error {
	return errors.New("Cannot invoke audio command. Must invoke a sub-command. See help for more info")
}

func onAudioPlayInvoked(ctx *commands.CommandContext) error {
	if len(ctx.Args) == 0 || len(ctx.Event.Mentions) == 0 {
		return errors.New("Need to provide an audio file name as the first argument, and all subsequent arguments should be target mentions")
	}

	// Check that the first argument isn't a mention (of the form <@!id>)
	audioFileName := ctx.Args[0]
	if audioFileName[0:3] == "<@!" {
		return errors.New("First argument needs to be a file name, not a mention")
	}
	// Make sure audio file exists before moving and adding to roles
	sanitisedFileName := strings.Replace(audioFileName, "\\", "", -1)
	sanitisedFileName = strings.Replace(sanitisedFileName, "/", "", -1)
	if !audio.Exists(audioFileName) {
		return errors.New("Could not find audio file " + sanitisedFileName)
	}

	message := ctx.Event
	textChannel, _ := ctx.Session.State.Channel(message.ChannelID)
	guildID := textChannel.GuildID
	guild, _ := ctx.Session.State.Guild(guildID)

	voiceChannelID := config.Config.Channel
	// Ensure the voice channel exists still
	_, err := ctx.Session.State.Channel(voiceChannelID)
	if err != nil {
		return errors.Wrap(err, "Could not find voice channel for guild.")
	}

	roleID := config.Config.Role
	// Ensure the role exists still
	_, err = ctx.Session.State.Role(guildID, roleID)
	if err != nil {
		return errors.Wrap(err, "Could not find locking role for guild.")
	}

	// Get their current channels to return them to
	originalChannels := make(map[string]string)
	for _, user := range message.Mentions {
		for _, voiceState := range guild.VoiceStates {
			if voiceState.UserID == user.ID {
				originalChannels[user.ID] = voiceState.ChannelID
			}
		}
	}

	// Defer moving users back before role removal
	defer func() {
		for _, user := range message.Mentions {
			ctx.Session.GuildMemberMove(guildID, user.ID, originalChannels[user.ID])
		}
	}()

	// Defer removing them from the locked roles
	defer func() {
		for _, user := range message.Mentions {
			ctx.Session.GuildMemberRoleRemove(guildID, user.ID, roleID)
		}
	}()

	// Add to locked role
	for _, user := range message.Mentions {
		ctx.Session.GuildMemberRoleAdd(guildID, user.ID, roleID)
	}

	// Move into the target channel
	for _, user := range message.Mentions {
		err = ctx.Session.GuildMemberMove(guildID, user.ID, voiceChannelID)
	}

	audio.PlaybackAudio(ctx.Args[0], guildID, voiceChannelID, &ctx.Session)

	return nil
}

func createAudioCommand() *commands.Command {
	audioParentCommand.AddChild(&audioPlayCommand)

	return &audioParentCommand
}
