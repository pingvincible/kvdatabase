package parser_test

import (
	"testing"

	"github.com/pingvincible/kvdatabase/internal/compute/parser"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) { //nolint: funlen // test code
	cases := []struct {
		name        string
		text        string
		wantCommand parser.Command
		wantError   error
	}{
		{
			name:        "empty command text",
			text:        "",
			wantCommand: parser.Command{},
			wantError:   parser.ErrInvalidCommand,
		},
		{
			name:        "Incorrect command name",
			text:        "INVALID",
			wantCommand: parser.Command{},
			wantError:   parser.ErrInvalidCommand,
		},
		{
			name:        "command name not in upper register",
			text:        "set",
			wantCommand: parser.Command{},
			wantError:   parser.ErrInvalidCommand,
		},
		{
			name: "SET correct command",
			text: "SET key value",
			wantCommand: parser.Command{
				Type:  parser.CommandSet,
				Key:   "key",
				Value: "value",
			},
			wantError: nil,
		},
		{
			name: "SET command with more than 2 arguments",
			text: "SET key value extraValue",
			wantCommand: parser.Command{
				Type:  parser.CommandSet,
				Key:   "key",
				Value: "value",
			},
			wantError: nil,
		},
		{
			name:        "SET command with no arguments",
			text:        "SET",
			wantCommand: parser.Command{},
			wantError:   parser.ErrNotEnoughArguments,
		},
		{
			name:        "SET command with one argument",
			text:        "SET key",
			wantCommand: parser.Command{},
			wantError:   parser.ErrNotEnoughArguments,
		},
		{
			name:        "SET command with empty arguments",
			text:        "SET   ",
			wantCommand: parser.Command{},
			wantError:   parser.ErrNotEnoughArguments,
		},
		{
			name:        "SET command with invalid key",
			text:        "SET invalid#key value",
			wantCommand: parser.Command{},
			wantError:   parser.ErrInvalidArgument,
		},
		{
			name:        "SET command with invalid value",
			text:        "SET key invalid& value",
			wantCommand: parser.Command{},
			wantError:   parser.ErrInvalidArgument,
		},
		{
			name: "GET correct command",
			text: "GET key",
			wantCommand: parser.Command{
				Type:  parser.CommandGet,
				Key:   "key",
				Value: "",
			},
			wantError: nil,
		},
		{
			name:        "GET command without arguments",
			text:        "GET",
			wantCommand: parser.Command{},
			wantError:   parser.ErrNotEnoughArguments,
		},
		{
			name:        "GET command with empty arguments",
			text:        "GET ",
			wantCommand: parser.Command{},
			wantError:   parser.ErrNotEnoughArguments,
		},
		{
			name:        "GET command with invalid key",
			text:        "GET invalid#key",
			wantCommand: parser.Command{},
			wantError:   parser.ErrInvalidArgument,
		},
		{
			name: "DEL correct command",
			text: "DEL key",
			wantCommand: parser.Command{
				Type:  parser.CommandDel,
				Key:   "key",
				Value: "",
			},
			wantError: nil,
		},
		{
			name:        "DEL command without arguments",
			text:        "DEL",
			wantCommand: parser.Command{},
			wantError:   parser.ErrNotEnoughArguments,
		},
		{
			name:        "DEL command with empty arguments",
			text:        "DEL ",
			wantCommand: parser.Command{},
			wantError:   parser.ErrNotEnoughArguments,
		},
		{
			name:        "DEL command with invalid key",
			text:        "DEL invalid#key",
			wantCommand: parser.Command{},
			wantError:   parser.ErrInvalidArgument,
		},
	}

	t.Parallel()

	for _, tc := range cases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			command, err := parser.Parse(testCase.text)
			assert.Equal(t, testCase.wantCommand, command)
			assert.Equal(t, testCase.wantError, err)
		})
	}
}
