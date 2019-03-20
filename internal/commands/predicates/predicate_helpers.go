package predicates

import (
	"fmt"
	"kingsgate/internal/commands"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// HasPermission finds if a member has a permission
func HasPermission(session *discordgo.Session, member *discordgo.Member, permission int) (hasPermission bool, err error) {
	roles, err := session.GuildRoles(member.GuildID)
	if err != nil {
		return false, err
	}

	// Iterate roles in guild
	for _, role := range roles {
		// If the permission int Bitwise And the permission gives the permission, means they have it
		if role.Permissions&permission == permission {
			// Check member's roles, see if they have it
			for _, roleID := range member.Roles {
				if role.ID == roleID {
					return true, nil
				}
			}
		}
	}
	err = errors.New(fmt.Sprintf("Could not find a role with the permission \"%d\" that the user is in", permission))
	return false, err
}

func getMember(ctx *commands.CommandContext, memberID string) (member *discordgo.Member, err error) {
	channel, err := ctx.Session.State.Channel(ctx.Event.ChannelID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to find the channel the message originated in")
	}

	guild, err := ctx.Session.State.Guild(channel.GuildID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to find the guild the message originated")
	}

	member, err = ctx.Session.GuildMember(guild.ID, memberID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to find the member to calculate permissions for")
	}

	// For some reason member.GuildID isn't set by GuildMember
	member.GuildID = guild.ID
	return member, nil
}
