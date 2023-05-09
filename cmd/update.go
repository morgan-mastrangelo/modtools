package cmd

import (
	"github.com/muesli/coral"
	"golang.org/x/mod/semver"
)

func init() {
	rootCmd.AddCommand(&coral.Command{
		Use:   "update",
		Short: "Update dependencies to the latest version",
		Long:  "Updates all direct dependencies to a newer version if available so that 'modtools check' passes",
		RunE: func(cmd *coral.Command, args []string) error {
			return updateDeps(getDirectOnly(cmd))
		},
	})
}

func updateDeps(directOnly bool) error {
	e, err := loadExceptions()
	if err != nil {
		return err
	}
	list, err := readDeps(true, directOnly)
	if err != nil {
		return err
	}
	for _, item := range list {
		if item.Update.Version != "" {
			if ex := e.Get(item.Path); ex != nil {
				if semver.Compare(item.Version, ex.MinVersion) >= 0 {
					continue
				}
			}
			_, err = runCommand("go", "get", "-d", item.Path+"@"+item.Update.Version)
			if err != nil {
				break
			}
		}
	}
	if err == nil {
		err = e.Save()
	}
	return err
}
