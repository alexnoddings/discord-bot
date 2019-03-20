package predicates

import (
	"fmt"
	"kingsgate/internal/commands"

	"github.com/pkg/errors"
)

// IsGuildAdminPredicate finds if the author of an event is an admin in the guild the event took place in
func IsGuildAdminPredicate(ctx *commands.CommandContext) error {
	member, err := getMember(ctx, ctx.Event.Author.ID)
	if err != nil {
		return err
	}

	isAdmin, err := HasPermission(&ctx.Session, member, 8)
	if !isAdmin {
		return errors.New("Only guild admins may use this command")
	}

	if err != nil {
		return err
	}

	return nil
}

// IsBotOwnerPredicate finds if the author of an event is the owner of the bot
func IsBotOwnerPredicate(ctx *commands.CommandContext) error {
	if ctx.Event.Author.ID != "128557356106645504" {
		return errors.New("Only the bot owner may use this command")
	}

	return nil
}

// GenerateAuthorHasPermissionPredicate finds if the author of an event has a given permission in the guild the event took place in
func GenerateAuthorHasPermissionPredicate(permission int) func(ctx *commands.CommandContext) error {
	return func(ctx *commands.CommandContext) error {
		member, err := getMember(ctx, ctx.Event.Author.ID)
		if err != nil {
			return errors.Wrap(err, "Failed to find the member to calculate permissions for")
		}

		hasPermission, err := HasPermission(&ctx.Session, member, permission)
		if err != nil {
			return errors.Wrap(err, "Failed to calculate if author has permission")
		}

		if hasPermission {
			return nil
		}
		return errors.New(fmt.Sprintf("User does not have the required permissions (%d)", permission))
	}
}

// GenerateBotHasPermissionPredicate finds if the bot has a given permission in the guild the event took place in
func GenerateBotHasPermissionPredicate(permission int) func(ctx *commands.CommandContext) error {
	return func(ctx *commands.CommandContext) error {
		self, err := getMember(ctx, ctx.Event.Author.ID)
		if err != nil {
			return errors.Wrap(err, "Failed to find the member to calculate permissions for")
		}

		hasPermission, err := HasPermission(&ctx.Session, self, permission)
		if err != nil {
			return errors.Wrap(err, "Failed to calculate if bot has permission")
		}

		if hasPermission {
			return nil
		}
		return errors.New(fmt.Sprintf("Bot does not have the required permissions (%d)", permission))
	}
}
