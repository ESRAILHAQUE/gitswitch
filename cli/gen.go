package cli

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ESRAILHAQUE/gitswitch/core/settings"
	"github.com/ESRAILHAQUE/gitswitch/ui"
	"github.com/ESRAILHAQUE/gitswitch/core/keygen"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a new SSH key and register a profile",
	RunE:  runGen,
}

func runGen(_ *cobra.Command, _ []string) error {
	cfg, err := settings.Load()
	if err != nil {
		return err
	}

	name, err := prompt.Ask("Full name", "")
	if err != nil {
		return err
	}
	email, err := prompt.Ask("Email", "")
	if err != nil {
		return err
	}
	key, err := prompt.Ask("Key (short alias, e.g. work, personal)", "")
	if err != nil {
		return err
	}
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	if cfg.FindByKey(key) != nil {
		return fmt.Errorf("profile %q already exists — run 'gitswitch del %s' first", key, key)
	}

	keyPath, err := keygen.GenerateKey(key, email)
	if err != nil {
		return err
	}
	if err := keygen.AddToSSHConfig(key, keyPath); err != nil {
		return fmt.Errorf("failed to update SSH config: %w", err)
	}

	p := settings.Profile{Key: key, Label: key, Name: name, Email: email, SSHKey: keyPath}
	if err := cfg.Add(p); err != nil {
		return err
	}
	if err := cfg.Save(); err != nil {
		return err
	}

	pubKey, err := os.ReadFile(keyPath + ".pub")
	if err != nil {
		return err
	}

	fmt.Printf("\nProfile %q created.\n", key)
	fmt.Printf("SSH host alias: git.%s → github.com\n\n", key)
	fmt.Printf("Add this public key to GitHub (Settings → SSH and GPG keys):\n\n%s\n", strings.TrimSpace(string(pubKey)))
	printSSHAddHint(keyPath)
	return nil
}

func printSSHAddHint(keyPath string) {
	if runtime.GOOS == "windows" {
		fmt.Printf("\nTo load the key (if using a passphrase or ssh-agent):\n")
		fmt.Printf("  Get-Service ssh-agent | Set-Service -StartupType Manual\n")
		fmt.Printf("  Start-Service ssh-agent\n")
		fmt.Printf("  ssh-add %s\n", keyPath)
		fmt.Printf("Open a new terminal if ssh-add is not found.\n")
		return
	}
	fmt.Printf("\nTo load the key in your current shell:\n  ssh-add %s\n", keyPath)
}
