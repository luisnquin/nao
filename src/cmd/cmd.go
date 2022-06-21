package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/src/core"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   core.AppName,
	Short: core.AppName + " is a tool to manage your notes",
	Long: `A tool to manage your notes or other types of files without
		worry about the path where it is, agile and safe if you want`,
	Run: func(cmd *cobra.Command, args []string) {
		/*
			password := []byte("oatasdadniasndasipdnipasndasodma")

			c, err := security.EncryptContent(password, []byte("as9d7}}as9d7as97d98sa7d987asdsa7das "))
			if err != nil {
				panic(err)
			}

			fmt.Println(c)
			fmt.Println(string(c))

			s, err := security.DecryptContent(password, c)
			if err != nil {
				panic(err)
			}

			fmt.Println(s)
			fmt.Println(string(s))
		*/

		switch length := len(args); {
		case length == 0:
			// TODO: if draftByDefaultDisabled ...
			mainCmd.Run(cmd, args)

		case length == 1:
			editCmd.Run(cmd, args)

		default:
			cmd.Usage()
		}
	},
	TraverseChildren: true,
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	root.AddCommand(newCmd, renderCmd, mergeCmd, lsCmd, mainCmd, editCmd)
}
