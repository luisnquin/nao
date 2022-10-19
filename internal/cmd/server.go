package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/luisnquin/nao/internal/api"
	"github.com/spf13/cobra"
)

type serverComp struct {
	cmd     *cobra.Command
	port    string
	quiet   bool
	verbose bool
}

func buildServer() serverComp {
	c := serverComp{
		cmd: &cobra.Command{
			Use:           "server",
			Short:         "Serve content via REST server",
			Args:          cobra.NoArgs,
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().StringVarP(&c.port, "port", "p", ":3000", "port to listen (e.g.: \"XXXX\")")
	c.cmd.Flags().BoolVarP(&c.quiet, "quiet", "q", false, "keep the server quiet except in case of an exception")
	c.cmd.Flags().BoolVarP(&c.verbose, "verbose", "v", false, "start the server with a logger")

	return c
}

func (s *serverComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		if len(s.port) > 0 && string([]rune(s.port)[0]) != ":" {
			s.port = ":" + s.port
		}

		if !s.quiet {
			color.New(color.FgHiGreen).Fprintln(os.Stdout, "Listening on http://localhost"+s.port+"\n")
		}

		return api.New(s.port, s.quiet, s.verbose).Start()
	}
}
