package main

import (
	"golang.rgm.io/dwtk/cmd"
)

func main() {
	cmd.RootCmd.AddCommand(
		cmd.DisableCmd,
		cmd.DumpCmd,
		cmd.FlashCmd,
		cmd.GDBServerCmd,
		cmd.InfoCmd,
		cmd.ResetCmd,
	)
	cmd.RootCmd.Execute()
}
