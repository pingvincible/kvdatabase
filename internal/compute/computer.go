package compute

import (
	"fmt"

	"github.com/pingvincible/kvdatabase/internal/compute/parser"
)

type Setter interface {
	Set(key, value string)
}

type Getter interface {
	Get(key string) string
}

type Deleter interface {
	Delete(key string)
}

type StorageInterface interface {
	Setter
	Getter
	Deleter
}

type Computer struct {
	storage StorageInterface
}

func NewComputer(storage StorageInterface) *Computer {
	return &Computer{storage: storage}
}

func (c *Computer) Compute(command parser.Command) string {
	result := ""

	switch command.Type {
	case parser.CommandSet:
		c.storage.Set(command.Key, command.Value)
	case parser.CommandGet:
		result = c.storage.Get(command.Key)
	case parser.CommandDel:
		c.storage.Delete(command.Key)
	}

	return result
}

func (c *Computer) Process(text string) (string, error) {
	const op = "compute.computer.Process"

	command, err := parser.Parse(text)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return c.Compute(command), nil
}
