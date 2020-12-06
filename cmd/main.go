package cmd

import (
	"io"
	"log"
	"os"

	"github.com/RageZBla/cf-stale-route-detector/commands"
	"github.com/RageZBla/cf-stale-route-detector/detector"
	"github.com/RageZBla/cf-stale-route-detector/diego"
	"github.com/RageZBla/cf-stale-route-detector/gorouter"
	"github.com/RageZBla/cf-stale-route-detector/presenters"

	"github.com/jessevdk/go-flags"
)

type options struct {
	Version bool `                             short:"v"  long:"version"                                                            description:"prints the om release version"`
}

func Main(sout io.Writer, serr io.Writer, version string, args []string) error {
	stderr := log.New(serr, "", 0)
	stdout := log.New(sout, "", 0)

	var global options
	parser := flags.NewParser(&global, flags.PassDoubleDash|flags.PassAfterNonOption)
	parser.Name = "cf-stale-route-detector"

	args, _ = parser.ParseArgs(args[1:])

	if global.Version {
		return commands.NewVersion(version, sout).Execute(nil)
	}

	if len(args) > 0 && args[0] == "help" {
		args[0] = "--help"
	}

	presenter := presenters.NewLoggerPresenter(stdout, stderr)

	gorouterParser := gorouter.NewGorouterTableParser()
	actualMapper := diego.NewActualLRPMapper()
	desiredMapper := diego.NewDesiredLRPMapper()

	detector := detector.NewStaleRouteDetector(gorouterParser, actualMapper, desiredMapper)

	_, err := parser.AddCommand(
		"detect",
		"detect stale routes",
		"This command detects stale route using gorouter routing table and Diego actual LRPs exports.",
		commands.NewDetect(detector, presenter),
	)
	if err != nil {
		return err
	}

	_, err = parser.AddCommand(
		"version",
		"prints the cf-stale-route-detector version",
		"This command prints the cf-stale-route-detector version number.",
		commands.NewVersion(version, sout),
	)
	if err != nil {
		return err
	}

	parser.Options |= flags.HelpFlag

	_, err = parser.ParseArgs(args)
	if err != nil {
		if e, ok := err.(*flags.Error); ok {
			switch e.Type {
			case flags.ErrHelp, flags.ErrCommandRequired:
				parser.WriteHelp(os.Stdout)
				return nil
			}
		}
	}
	return err
}
