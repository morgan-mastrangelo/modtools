package cmd

import (
	"errors"
	"fmt"

	"github.com/muesli/coral"
	"golang.org/x/mod/semver"
)

func init() {
	rootCmd.AddCommand(&coral.Command{
		Use:   "check",
		Short: "Check for out-of-date dependencies",
		Long:  "Scans through direct and indirect dependencies to check if a newer version is available. Exceptions can be set by 'modtools freeze'",
		RunE: func(cmd *coral.Command, args []string) error {
			res, err := checkDeps(getDirectOnly(cmd))
			if err != nil {
				return err
			}
			if len(res) > 0 {
				fmt.Printf("Some dependencies are out-of-date. Please upgrade by running 'modtools update' or the following commands:\n\n")
				for _, item := range res {
					fmt.Printf("go get %s@%s\n", item.Path, item.Version)
				}
				fmt.Println()
				return errors.New("check has failed")
			}
			return nil
		},
	})
}

func checkDeps(directOnly bool) (updates []Update, err error) {
	e, err := loadExceptions()
	if err != nil {
		return nil, err
	}
	list, err := readDeps(true, directOnly)
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		if ex := e.Get(item.Path); ex != nil {
			if semver.Compare(item.Version, ex.MinVersion) < 0 {
				updates = append(updates, Update{Path: item.Path, Version: ex.MinVersion})
				fmt.Printf("Frozen module %q should be at least version %q\n", item.Path, ex.MinVersion)
			}
			continue
		}
		if item.Update.Version != "" {
			updates = append(updates, item.Update)
		}
	}
	return
}
