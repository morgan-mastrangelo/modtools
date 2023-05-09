package cmd

import (
	"github.com/muesli/coral"
)

func init() {
	rootCmd.AddCommand(&coral.Command{
		Use:   "thaw modpath",
		Short: "Un-freeze a dependency",
		Long:  "Removes the specified module path from the list of exceptions.",
		Args:  coral.ExactArgs(1),
		RunE: func(cmd *coral.Command, args []string) error {
			return thaw(args[0])
		},
	})
}

func thaw(p string) error {
	e, err := loadExceptions()
	if err != nil {
		return err
	}
	e.Remove(p)
	return e.Save()
}
