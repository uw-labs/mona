package command

import (
	"github.com/uw-labs/mona/internal/app"
	"github.com/uw-labs/mona/internal/config"
	"github.com/uw-labs/mona/internal/git"
)

type Config struct {
	Apps     []*app.App
	Diff     git.GoDiff
	Project  *config.Project
	FailFast bool
}

func NewConfig(wd, branch string) (cfg Config, err error) {
	cfg.Project, err = config.LoadProject(wd, branch)
	if err != nil {
		return cfg, err
	}

	cfg.Apps, err = app.FindApps("./", cfg.Project.Mod)
	if err != nil {
		return cfg, err
	}

	cfg.Diff, err = git.GetGoDiff(cfg.Project.Mod, cfg.Project.Branch)

	return cfg, err
}
