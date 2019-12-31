package cmd

import (
	"fmt"

	"github.com/dwtk/dwtk/debugwire"
	"github.com/dwtk/dwtk/internal/logger"
	"github.com/dwtk/dwtk/internal/version"
	"github.com/spf13/cobra"
)

var (
	dw      *debugwire.DebugWIRE
	noReset bool

	serialPort string
	baudrate   uint32
	debug      bool
)

func init() {
	RootCmd.PersistentFlags().StringVarP(
		&serialPort,
		"serial-port",
		"s",
		"",
		"serial port device (e.g. /dev/ttyUSB0) (Default: detect)",
	)
	RootCmd.PersistentFlags().Uint32VarP(
		&baudrate,
		"baudrate",
		"b",
		0,
		"serial port baudrate (e.g. 62500) (Default: detect)",
	)
	RootCmd.PersistentFlags().BoolVarP(
		&debug,
		"debug",
		"d",
		false,
		"enable debugging messages",
	)
	RootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
}

var RootCmd = &cobra.Command{
	Use:          "dwtk",
	Short:        "debugWIRE toolkit for AVR microcontrollers",
	Long:         "debugWIRE toolkit for AVR microcontrollers",
	Version:      version.Version,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			logger.EnableDebug()
		}

		var err error
		dw, err = debugwire.New(serialPort, baudrate)
		if err != nil {
			return err
		}

		if dw.MCU == nil {
			return fmt.Errorf("Failed to detect MCU")
		}

		noReset = false

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if dw != nil {
			defer dw.Close()
			if !noReset {
				if err := dw.Reset(); err != nil {
					return err
				}
				if err := dw.Go(); err != nil {
					return err
				}
			}
		}

		return nil
	},
}