package command

import (
	"fmt"
	"strings"
)

type Command struct {
	command []string
	desc    string
}

func NewCommand(desc string, cms ...string) Command {
	return Command{
		command: cms,
		desc:    desc,
	}
}

type Commands []Command

func NewCommands(cms ...Command) Commands {
	var c Commands
	c = append(c, cms...)
	return c
}

func (c Commands) String() string {
	var builder strings.Builder
	for _, command := range c {
		for _, oneCommand := range command.command {
			builder.WriteString(fmt.Sprintf("%s ", oneCommand))
		}
		builder.WriteString(fmt.Sprintf(": %s\n", command.desc))
	}
	return builder.String()
}
