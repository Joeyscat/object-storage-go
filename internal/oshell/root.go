package oshell

import (
	"github.com/spf13/cobra"
)

var (
	// oshell 版本信息， oshell -v
	VersionFlag bool
)

// cobra root cmd, all other commands is children or subchildren of this root cmd
var RootCmd = &cobra.Command{
	Use:     "oshell",
	Short:   "A commandline tool for managing your object storage",
	Version: version,
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&VersionFlag, "version", "v", false, "show version")
}
