package commands

import (
	"errors"
	"fmt"
	"strings"
)

type Command struct {
	// parent command
	parent *Command

	// Name of the command
	Name string
	
	// Aliases it can be invoked by
	Aliases []string

	// Predicates that need to return a nil error for the Command to be invoked
	RequiredPredicates []func(ctx *CommandContext) error

	// Function that is invoked when the command is
	OnInvoked func(ctx *CommandContext) error

	// Child functions
	children []*Command
}

// Checks if a command should be invoked by a certain command
func (command *Command) Matches(arg string) bool {
	arg = strings.ToLower(arg)
	if strings.ToLower(command.Name) == arg {
		return true
	}
	for _, alias := range command.Aliases {
		if strings.ToLower(alias) == arg {
			return true
		}
	}
	return false
}

func (command *Command) CanInvoke(ctx *CommandContext) (bool, error) {
	if command.parent != nil {
		parentCanInvoke, err := command.parent.CanInvoke(ctx)
		if !parentCanInvoke || err != nil {
			return false, err
		}
	}
	for _, predicate := range command.RequiredPredicates {
		err := predicate(ctx)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (command *Command) Invoke(ctx *CommandContext) error {
	canInvoke, err := command.CanInvoke(ctx)
	if !canInvoke || err != nil {
		return err
	}
	err = command.OnInvoked(ctx)
	return err
}

func (command *Command) AddChild(child *Command) error {
	if child.parent != nil {
		return errors.New(fmt.Sprintf("Child command %s already has a parent, %s", child.Name, command.Name))
	}
	child.parent = command
	command.children = append(command.children, child)
	return nil
}
