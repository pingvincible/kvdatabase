package parser

import "fmt"

type CommandType string

const (
	CommandSet CommandType = "SET"
	CommandGet CommandType = "GET"
	CommandDel CommandType = "DEL"

	CommandGetDelArgsCount = 1
	CommandSetArgsCount    = 2
)

type Command struct {
	Type  CommandType
	Key   string
	Value string
}

func (c *Command) String() string {
	return fmt.Sprintf("Type: %s, %s=%s", c.Type, c.Key, c.Value)
}

func ArgsInCommand() map[CommandType]int {
	return map[CommandType]int{
		CommandSet: CommandSetArgsCount,
		CommandGet: CommandGetDelArgsCount,
		CommandDel: CommandGetDelArgsCount,
	}
}
