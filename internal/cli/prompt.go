package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// isStdinTerminal reports whether stdin is connected to a terminal, so we
// can skip interactive prompting in CI or when output is piped.
func isStdinTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// prompt reads a line from stdin. If the line is empty, def is returned.
// A trailing newline is stripped.
func prompt(reader *bufio.Reader, label, def string) (string, error) {
	if def != "" {
		fmt.Printf("%s [%s]: ", label, def)
	} else {
		fmt.Printf("%s: ", label)
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return def, nil
	}
	return line, nil
}

// promptSelect presents numbered options and returns the chosen one.
// defIndex is the zero-based index of the default (or -1 for none).
func promptSelect(reader *bufio.Reader, label string, options []string, defIndex int) (string, error) {
	fmt.Printf("%s:\n", label)
	for i, opt := range options {
		marker := "  "
		if defIndex == i {
			marker = "* "
		}
		fmt.Printf("%s%d) %s\n", marker, i+1, opt)
	}

	defStr := ""
	if defIndex >= 0 {
		defStr = options[defIndex]
	}

	for {
		line, err := prompt(reader, "Choose", defStr)
		if err != nil {
			return "", err
		}

		if line == "" && defIndex >= 0 {
			return defStr, nil
		}

		var idx int
		if _, err := fmt.Sscanf(line, "%d", &idx); err != nil {
			fmt.Println("Please enter a number.")
			continue
		}
		idx-- // convert from 1-based to 0-based
		if idx < 0 || idx >= len(options) {
			fmt.Printf("Please enter a number between 1 and %d.\n", len(options))
			continue
		}
		return options[idx], nil
	}
}

// promptBool asks a yes/no question. def is the value returned on empty input.
func promptBool(reader *bufio.Reader, label string, def bool) (bool, error) {
	yesno := "Y/n"
	if !def {
		yesno = "y/N"
	}
	for {
		line, err := prompt(reader, label+" ("+yesno+")", "")
		if err != nil {
			return false, err
		}
		line = strings.ToLower(strings.TrimSpace(line))
		if line == "" {
			return def, nil
		}
		switch line {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			fmt.Println("Please answer y or n.")
		}
	}
}

// promptNonEmpty reads until a non-empty line is entered, returning def if
// the user just presses enter.
func promptNonEmpty(reader *bufio.Reader, label, def string) (string, error) {
	for {
		line, err := prompt(reader, label, def)
		if err != nil {
			return "", err
		}
		if line != "" {
			return line, nil
		}
		fmt.Println("Value is required.")
	}
}

// fieldTypeOptions are the supported field types for a resource.
var fieldTypeOptions = []string{"string", "int", "float", "bool", "time", "uuid"}

// communicationOptions are the supported communication protocols.
var communicationOptions = []string{"grpc"}

// architectureOptions are the supported code layouts.
var architectureOptions = []string{"vertical", "horizontal"}

// databaseDriverOptions are the supported database drivers.
var databaseDriverOptions = []string{"postgres", "mysql", "sqlite"}

// ciProviderOptions are the supported CI providers.
var ciProviderOptions = []string{"github", "gitlab"}

// cloudProviderOptions are the supported cloud providers.
var cloudProviderOptions = []string{"none", "aws"}
