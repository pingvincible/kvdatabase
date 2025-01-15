package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type CommandType string

const (
	CommandSet CommandType = "SET"
	CommandGet CommandType = "GET"
	CommandDel CommandType = "DEL"

	CommandSetArgsNumber = 2
)

type Command struct {
	Type  CommandType
	Key   string
	Value string
}

func (c *Command) String() string {
	return fmt.Sprintf("Type: %s, %s=%s", c.Type, c.Key, c.Value)
}

var (
	ErrUnknownCommand     = errors.New("unknown command")
	ErrNotEnoughArguments = errors.New("not enough arguments")
	ErrInvalidArgument    = errors.New("invalid argument")
)

func Parse(command string) (Command, error) {
	if len(command) == 0 {
		return Command{}, ErrUnknownCommand
	}

	args := strings.Split(command, ` `)
	if len(args) == 1 {
		return Command{}, ErrNotEnoughArguments
	}

	switch args[0] {
	case string(CommandSet):
		return parseCommandSet(args[1:])
	case string(CommandGet):
		return parseCommandGet(args[1:])
	case string(CommandDel):
		return parseCommandDel(args[1:])
	default:
		return Command{}, ErrUnknownCommand
	}
}

func parseCommandSet(args []string) (Command, error) {
	if len(args) < CommandSetArgsNumber {
		return Command{}, ErrNotEnoughArguments
	}

	err := validate(args)
	if err != nil {
		return Command{}, err
	}

	return Command{
		Type:  CommandSet,
		Key:   args[0],
		Value: args[1],
	}, nil
}

func parseCommandGet(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, ErrNotEnoughArguments
	}

	err := validate(args)
	if err != nil {
		return Command{}, err
	}

	return Command{
		Type:  CommandGet,
		Key:   args[0],
		Value: "",
	}, nil
}

func parseCommandDel(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, ErrNotEnoughArguments
	}

	err := validate(args)
	if err != nil {
		return Command{}, err
	}

	return Command{
		Type:  CommandDel,
		Key:   args[0],
		Value: "",
	}, nil
}

func validate(args []string) error {
	for _, arg := range args {
		err := validateArgument(arg)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateArgument(arg string) error {
	r := regexp.MustCompile("^[a-zA-Z0-9_/*]+$")

	matched := r.MatchString(arg)
	if !matched {
		return ErrInvalidArgument
	}

	return nil
}
