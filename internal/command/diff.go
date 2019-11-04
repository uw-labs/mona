package command

// DiffSummary lists all the app names that changed.
type DiffSummary struct {
	All         []string
	ChangedApps []string
}

// Diff outputs the names of all apps where changes are detected.
func Diff(cfg Config) (summary DiffSummary, err error) {

	for _, appInfo := range cfg.Apps {
		summary.All = append(summary.All, appInfo.Name)
	}
	for _, appInfo := range getChangedApps(cfg) {
		summary.ChangedApps = append(summary.ChangedApps, appInfo.Name)
	}

	return summary, nil
}
