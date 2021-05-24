package config

import (
	"strings"

	"github.com/SmallTianTian/fresh-go/config"
)

func GetModule(c *config.Config) string {
	remote := c.Project.Remote
	owner := c.Project.Owner
	name := c.Project.Name

	if name == "" {
		panic("No project name.")
	}

	if owner != "" && remote == "" {
		panic("No remote must no owner")
	}

	var mods []string
	for _, v := range []string{remote, owner, name} {
		if v != "" {
			mods = append(mods, v)
		}
	}

	return strings.Join(mods, "/")
}
