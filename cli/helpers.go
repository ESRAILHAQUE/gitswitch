package cli

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ESRAILHAQUE/gitswitch/core/settings"
	"github.com/ESRAILHAQUE/gitswitch/ui"
)

var httpsPattern = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+?)(?:\.git)?$`)

func resolveRemoteURL(raw string, cfg *settings.Config, promptLabel string) (string, bool, error) {
	switch {
	case strings.HasPrefix(raw, "https://github.com/") || strings.HasPrefix(raw, "http://github.com/"):
		path := strings.TrimPrefix(strings.TrimPrefix(raw, "https://"), "http://")
		path = strings.TrimPrefix(path, "github.com/")
		path = strings.TrimSuffix(path, ".git")

		if len(cfg.Profiles) == 0 {
			return "", false, fmt.Errorf("no profiles found — run 'gitswitch gen' to create one")
		}
		idx, err := prompt.Select(promptLabel, profileLabels(cfg.Profiles))
		if err != nil {
			return "", false, err
		}
		return fmt.Sprintf("git@git.%s:%s", cfg.Profiles[idx].Key, path), true, nil

	case strings.Contains(raw, ":") && !strings.Contains(raw, "://") && !strings.HasPrefix(raw, "git@"):
		parts := strings.SplitN(raw, ":", 2)
		key, path := parts[0], parts[1]
		if cfg.FindByKey(key) == nil {
			return "", false, fmt.Errorf("unknown profile key %q — run 'gitswitch list' to see available profiles", key)
		}
		path = strings.TrimSuffix(path, ".git")
		return fmt.Sprintf("git@git.%s:%s", key, path), true, nil

	default:
		return raw, false, nil
	}
}

func profileLabels(profiles []settings.Profile) []string {
	labels := make([]string, len(profiles))
	for i, p := range profiles {
		labels[i] = p.Label
	}
	return labels
}
