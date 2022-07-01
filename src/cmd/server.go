package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/luisnquin/nao/src/api"
	"github.com/spf13/cobra"
)

type serverComp struct {
	cmd  *cobra.Command
	port string
}

var server = buildServer()

func buildServer() serverComp {
	c := serverComp{
		cmd: &cobra.Command{
			Use:           "server",
			Short:         "Serve content via REST server",
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().StringVarP(&c.port, "port", "p", ":3000", "Port to listen (e.g.: \"XXXX\")")

	return c
}

func (s *serverComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		if len(s.port) > 0 && string([]rune(s.port)[0]) != ":" {
			s.port = ":" + s.port
		}

		color.New(color.FgHiGreen).Fprintln(os.Stdout, "Listening on http://localhost"+s.port+"\n")

		return api.New().Start(s.port)
	}
}
