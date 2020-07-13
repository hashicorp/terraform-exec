package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/logutils"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"github.com/mitchellh/cli"
)

func main() {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	ui := &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,
		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	exitStatus := run(ui, os.Args[1:])

	os.Exit(exitStatus)
}

func help() string {
	return `Usage: tfinstall [--dir=DIR] VERSION

  Downloads, verifies, and installs a terraform binary of version
  VERSION from releases.hashicorp.com. VERSION must be a valid version under
  semantic versioning, or "latest".

  If a binary is successfully installed, its path will be printed to stdout.

  Unless --dir is given, the default system temporary directory will be used.

Options:
  --dir          Directory into which to install the terraform binary. The
                 directory must exist.

Examples:
  tfinstall 0.12.28
  tfinstall latest
  tfinstall 0.13.0-beta3
  tfinstall --dir=/home/kmoe/bin 0.12.28
`
}

func run(ui cli.Ui, args []string) int {
	ctx := context.Background()

	args = os.Args[1:]
	flags := flag.NewFlagSet("", flag.ExitOnError)
	var tfDir string
	flags.StringVar(&tfDir, "dir", "", "Local directory into which to install terraform")

	err := flags.Parse(args)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	if flags.NArg() != 1 {
		ui.Error("Please specify VERSION")
		ui.Output(help())
		return 127
	}

	tfVersion := flags.Args()[0]

	if tfDir == "" {
		tfDir, err = ioutil.TempDir("", "tfinstall")
		if err != nil {
			ui.Error(err.Error())
			return 1
		}
	}

	var findArgs []tfinstall.ExecPathFinder

	if tfVersion == "latest" {
		findArgs = append(findArgs, tfinstall.LatestVersion(tfDir, false))
	} else {
		findArgs = append(findArgs, tfinstall.ExactVersion(tfVersion, tfDir))
	}

	tfPath, err := tfinstall.Find(ctx, findArgs...)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	ui.Output(tfPath)
	return 0
}
