package parser

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidCommand     = errors.New("invalid command")
	ErrNotEnoughArguments = errors.New("not enough arguments")
	ErrInvalidArgument    = errors.New("invalid argument")
)

func Parse(commandText string) (Command, error) {
	commandType, validatedArgs, err := validate(commandText)
	if err != nil {
		return Command{}, err
	}

	command := Command{Type: commandType, Key: validatedArgs[0]}
	if commandType == CommandSet {
		command.Value = validatedArgs[1]
	}

	return command, nil
}

func validate(commandText string) (CommandType, []string, error) {
	args := strings.Fields(commandText)
	if len(args) == 0 {
		return "", nil, ErrInvalidCommand
	}

	commandType, argsCount, err := validateCommand(CommandType(args[0]))
	if err != nil {
		return "", nil, err
	}

	err = validateArgs(args[1:], argsCount)
	if err != nil {
		return "", nil, err
	}

	return commandType, args[1:], nil
}

func validateCommand(commandType CommandType) (CommandType, int, error) {
	argsCount, ok := ArgsInCommand()[commandType]
	if !ok {
		return "", 0, ErrInvalidCommand
	}

	return commandType, argsCount, nil
}

func validateArgs(args []string, argsCount int) error {
	if len(args) < argsCount {
		return ErrNotEnoughArguments
	}

	for _, arg := range args {
		err := validateArg(arg)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateArg(arg string) error {
	r := regexp.MustCompile("^[a-zA-Z0-9_/*]+$")

	matched := r.MatchString(arg)
	if !matched {
		return ErrInvalidArgument
	}

	return nil
}
