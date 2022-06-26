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
	q rune = 113
	Q rune = 81
)

var exposeCmd = &cobra.Command{ // TODO: add support for fswatch
	Use:   "expose",
	Short: "Exposes all the sets in a directory",
	Run: func(cmd *cobra.Command, args []string) {
		views := data.New().ListSets()

		// timeout, _ := cmd.Flags().GetInt("timeout")
		// watch, _ := cmd.Flags().GetBool("watch")
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

			if untree || v.Type == constants.TypeMain {
				f, err = os.Create(config.App.Paths.CacheDir + "/" + v.Tag + "-" + v.Key[:5])
			} else if v.Extension != "" {
				f, err = os.Create(config.App.Paths.CacheDir + "/" + v.Type + "/" + v.Tag + "-" + v.Key[:5] + "." + v.Extension)
			} else {
				f, err = os.Create(config.App.Paths.CacheDir + "/" + v.Type + "/" + v.Tag + "-" + v.Key[:5])
			}

			cobra.CheckErr(err)

			_, err = f.WriteString(v.Content)
			cobra.CheckErr(err)

			_ = f.Close()
		}

		fmt.Fprintln(os.Stdout, "Files exposed on "+config.App.Paths.CacheDir)

		if detach {
			os.Exit(0)
		}

		cobra.CheckErr(keyboard.Open())

		defer keyboard.Close()

		fmt.Fprintln(os.Stdout, "Click Q or Ctrl+C to exit")

		for {
			char, key, err := keyboard.GetKey()
			cobra.CheckErr(err)

			if key == keyboard.KeyCtrlC || char == Q || char == q {
				break
			}
		}

		_ = os.RemoveAll(config.App.Paths.CacheDir)
		_ = os.MkdirAll(config.App.Paths.CacheDir, os.ModePerm)
	},
}

func init() {
	exposeCmd.Flags().BoolP("untree", "u", false, "disable default tree file organization depending on types")
	exposeCmd.Flags().BoolP("detach", "d", false, "")
	// exposeCmd.Flags().BoolP("watch", "w", false, "")
	//	exposeCmd.Flags().IntP("timeout", "t", 0, "set a time limit expresed in seconds for the exposition of files")
}
