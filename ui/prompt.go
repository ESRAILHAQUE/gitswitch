package prompt

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

// Ask prompts the user for text input.
func Ask(label, defaultVal string) (string, error) {
	var input string
	err := huh.NewInput().
		Title(label).
		Value(&input).
		Placeholder(defaultVal).
		Run()
	if err != nil {
		return "", err
	}
	if input == "" {
		return defaultVal, nil
	}
	return input, nil
}

// Select displays an interactive arrow-key menu to select an option.
func Select(label string, options []string) (int, error) {
	if len(options) == 0 {
		return 0, fmt.Errorf("no options available")
	}

	// We create a slice of Options that hold both the string (label) and the index (value)
	var huhOptions []huh.Option[int]
	for i, opt := range options {
		huhOptions = append(huhOptions, huh.NewOption(opt, i))
	}

	var selectedIndex int
	err := huh.NewSelect[int]().
		Title(label).
		Options(huhOptions...).
		Value(&selectedIndex).
		Run()
	if err != nil {
		return 0, err
	}

	return selectedIndex, nil
}

// Confirm asks a boolean Yes/No question.
func Confirm(label string) (bool, error) {
	var confirm bool
	err := huh.NewConfirm().
		Title(label).
		Value(&confirm).
		Run()
	if err != nil {
		return false, err
	}
	return confirm, nil
}
