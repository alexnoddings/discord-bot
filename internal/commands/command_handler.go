package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"strings"
)

const (
	indentString    = "   "
	pipeString      = " │ "
	intersectString = " ├─"
	cornerString    = " └─"
)

type CommandHandler struct {
	// What messages must begin with for the handler to parse them
	CommandPrefix string

	// Using a "root command" allows us to recursively find commands easier
	rootCommand Command
}

func NewCommandHandler(commandPrefix string) CommandHandler {
	return CommandHandler{
		CommandPrefix: commandPrefix,
		rootCommand:   Command{Name: "Root"},
	}
}

func (handler *CommandHandler) AddRootCommand(command *Command) error {
	// Can't add commands unless they are a root command (i.e. have no parent)
	if command.parent != nil {
		return errors.New(fmt.Sprintf("Command %s already has a parent (%s)", command.Name, command.parent.Name))
	}

	// handler.commands = append(handler.commands, *command)
	handler.rootCommand.AddChild(command)
	return nil
}

func (handler *CommandHandler) OnMessageCreated(session *discordgo.Session, messageCreated *discordgo.MessageCreate) {
	// Ignore messages from bots
	if messageCreated.Author.Bot {
		return
	}

	content := messageCreated.Content
	// Don't parse messages that don't begin with the command prefix
	if !strings.HasPrefix(content, handler.CommandPrefix) {
		return
	}

	prefixLength := len(handler.CommandPrefix)
	// Generate arguments from after the prefix
	args := splitArguments(content[prefixLength:])

	// Search all parent commands, see if the first arg matches it
	commandToInvoke, argsToInvokeWith := subCommandFinder(&handler.rootCommand, args)
	if commandToInvoke == &handler.rootCommand {
		session.MessageReactionAdd(messageCreated.ChannelID, messageCreated.ID, "\u2753")
		return
	}

	ctx := CommandContext{Session: *session, Event: *messageCreated, Args: argsToInvokeWith}
	err := commandToInvoke.Invoke(&ctx)

	if err != nil {
		session.ChannelMessageSend(messageCreated.ChannelID, fmt.Sprintf("Error occured while executing command:\n```%s```", err.Error()))
	}
}

func (handler *CommandHandler) GenerateTreeView() string {
	return handler.rootCommand.GenerateTreeView("", true)
}

func (command *Command) GenerateTreeView(indent string, last bool) (out string) {
	out += indent
	if last {
		indent += indentString
		out += cornerString
	} else {
		indent += pipeString
		out += intersectString
	}
	out += command.Name + " [" + strings.Join(command.Aliases, ", ") + "]" + "\n"

	for i, child := range command.children {
		out += child.GenerateTreeView(indent, i == len(command.children)-1)
	}
	return out
}

func subCommandFinder(parent *Command, args []string) (*Command, []string) {
	if len(args) > 0 {
		for _, child := range parent.children {
			if child.Matches(args[0]) {
				return subCommandFinder(child, args[1:])
			}
		}
	}
	return parent, args
}

func splitArguments(message string) []string {
	possibleArgs := strings.Split(message, " ")
	var actualArgs []string
	for _, arg := range possibleArgs {
		if arg != "" {
			actualArgs = append(actualArgs, arg)
		}
	}
	return actualArgs
}
