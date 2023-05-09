package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/muesli/coral"
)

const (
	defaultFreezeDays = 7
)

func init() {
	rootCmd.AddCommand(&coral.Command{
		Use:   "freeze modpath [days]",
		Short: "Freeze a dependency",
		Long: "Adds the specified module path to the list of exceptions. The check and update commands will ignore this " +
			"module for the specified number of days\n(defaults to " + strconv.Itoa(defaultFreezeDays) +
			").\nDon't forget to add " + frozendepsFilename + " to the repository.",
		Args: coral.RangeArgs(1, 2),
		RunE: func(cmd *coral.Command, args []string) error {
			var days int
			if len(args) > 1 {
				var err error
				days, err = strconv.Atoi(args[1])
				if err != nil {
					return err
				}
			} else {
				days = defaultFreezeDays
			}
			until := time.Now().Add(time.Duration(days) * 24 * time.Hour).Truncate(time.Second)
			return freeze(args[0], until)
		},
	})
}

func freeze(p string, until time.Time) error {
	e, err := loadExceptions()
	if err != nil {
		return err
	}
	deps, err := readDeps(false, false)
	if err != nil {
		return err
	}
	curVersion := ""
	for _, item := range deps {
		if item.Path == p {
			curVersion = item.Version
			break
		}
	}
	if curVersion == "" {
		return fmt.Errorf("module %q is not a dependency", p)
	}
	e.Add(&Exception{
		Path:       p,
		MinVersion: curVersion,
		ValidUntil: until,
	})
	wasNew := e.IsNew()
	err = e.Save()
	if err == nil {
		if wasNew {
			fmt.Println("Don't forget to add " + frozendepsFilename + " to the repository.")
		}
	}
	return err
}
