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

type exposeComp struct {
	cmd     *cobra.Command
	detach  bool
	untree  bool
	watch   bool
	timeout int
}

var expose = buildExpose()

func buildExpose() exposeComp {
	c := exposeComp{
		cmd: &cobra.Command{
			Use:           "expose",
			Short:         "Exposes all the sets in a directory",
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().BoolVarP(&c.detach, "detach", "d", false, "")
	c.cmd.Flags().BoolVarP(&c.untree, "untree", "u", false, "disable default tree file organization depending on types")
	// c.cmd.Flags().IntVarP(&c.timeout, "timeout", "t", 0, "set a time limit expresed in seconds for the exposition of files")
	// c.cmd.Flags().BoolVarP(&c.watch, "watch", "w", false, "")

	return c
}

func (e *exposeComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		views := data.New().ListSets()

		_ = os.MkdirAll(config.App.Paths.CacheDir, os.ModePerm)

		if !e.untree {
			_ = os.MkdirAll(config.App.Paths.CacheDir+"/"+constants.TypeDefault, os.ModePerm)
			_ = os.MkdirAll(config.App.Paths.CacheDir+"/"+constants.TypeImported, os.ModePerm)
			_ = os.MkdirAll(config.App.Paths.CacheDir+"/"+constants.TypeMerged, os.ModePerm)
		}

		var err error

		for _, v := range views {
			var f *os.File

			if e.untree || v.Type == constants.TypeMain {
				f, err = os.Create(config.App.Paths.CacheDir + "/" + v.Tag + "-" + v.Key[:5])
			} else if v.Extension != "" {
				f, err = os.Create(config.App.Paths.CacheDir + "/" + v.Type + "/" + v.Tag + "-" + v.Key[:5] + "." + v.Extension)
			} else {
				f, err = os.Create(config.App.Paths.CacheDir + "/" + v.Type + "/" + v.Tag + "-" + v.Key[:5])
			}

			if err != nil {
				return err
			}

			_, err = f.WriteString(v.Content)
			if err != nil {
				return err
			}

			_ = f.Close()
		}

		fmt.Fprintln(os.Stdout, "Files exposed on "+config.App.Paths.CacheDir)

		if e.detach {
			return nil
		}

		if err = keyboard.Open(); err != nil {
			return err
		}

		defer keyboard.Close()

		fmt.Fprintln(os.Stdout, "Click Q or Ctrl+C to exit")

		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				return err
			}

			if key == keyboard.KeyCtrlC || char == Q || char == q {
				break
			}
		}

		_ = os.RemoveAll(config.App.Paths.CacheDir)
		_ = os.MkdirAll(config.App.Paths.CacheDir, os.ModePerm)

		return nil
	}
}
