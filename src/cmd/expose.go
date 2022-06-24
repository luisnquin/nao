package cmd

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/spf13/cobra"
)

const (
	keyq keyboard.Key = 113
	keyQ keyboard.Key = 81
)

var exposeCmd = &cobra.Command{
	Use:   "expose",
	Short: "Exposes all the sets in a directory",
	Long:  "...",
	Run: func(cmd *cobra.Command, args []string) {
		views := data.NewUserBox().ListSets()

		detach, _ := cmd.Flags().GetBool("detach")
		untree, _ := cmd.Flags().GetBool("untree")

		_ = os.MkdirAll(config.App.Paths.CacheDir, os.ModePerm)

		if !untree {
			_ = os.MkdirAll(config.App.Paths.CacheDir+"/"+constants.TypeDefault, os.ModePerm)
			_ = os.MkdirAll(config.App.Paths.CacheDir+"/"+constants.TypeImported, os.ModePerm)
			_ = os.MkdirAll(config.App.Paths.CacheDir+"/"+constants.TypeMerged, os.ModePerm)
		}

		var err error

		for _, v := range views {
			var f *os.File

			if !untree {
				f, err = os.Create(config.App.Paths.CacheDir + "/" + v.Type + "/" + v.Tag + "-" + v.Key[:5]) // Pending: sufix
			} else {
				f, err = os.Create(config.App.Paths.CacheDir + "/" + v.Tag + "-" + v.Key[:5]) // Pending: sufix
			}

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			_, err = f.WriteString(v.Content)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			_ = f.Close()
		}

		fmt.Fprintln(os.Stdout, "Files exposed on "+config.App.Paths.CacheDir)

		if detach {
			os.Exit(0)
		}

		if err = keyboard.Open(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		defer keyboard.Close()

		fmt.Fprintln(os.Stdout, "Click Q or Ctrl+C to exit")

		for {
			_, key, err := keyboard.GetKey()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if key == keyboard.KeyCtrlC || key == keyQ || key == keyq {
				_ = os.RemoveAll(config.App.Paths.CacheDir)
				_ = os.MkdirAll(config.App.Paths.CacheDir, os.ModePerm)

				break
			}
		}
	},
}

func init() {
	exposeCmd.Flags().Bool("untree", false, constants.AppName+" expose --untree")
	exposeCmd.Flags().Bool("detach", false, constants.AppName+" expose --detach")
}
