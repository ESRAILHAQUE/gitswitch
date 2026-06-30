package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ESRAILHAQUE/gitswitch/core/settings"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	RunE:  runList,
}

func runList(_ *cobra.Command, _ []string) error {
	cfg, err := settings.Load()
	if err != nil {
		return err
	}
	if len(cfg.Profiles) == 0 {
		fmt.Println("No profiles found. Run 'gitswitch gen' to create one.")
		return nil
	}
	fmt.Printf("%-12s  %-28s  %s\n", "KEY", "EMAIL", "SSH KEY")
	fmt.Printf("%-12s  %-28s  %s\n", "---", "-----", "-------")
	for _, p := range cfg.Profiles {
		fmt.Printf("%-12s  %-28s  %s\n", p.Key, p.Email, p.SSHKey)
	}
	return nil
}
