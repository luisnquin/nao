package cmd

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/andreaskoch/go-fswatch"
	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/luisnquin/nao/internal/config"
	"github.com/luisnquin/nao/internal/store"
	"github.com/spf13/cobra"
)

const (
	q rune = 113
	Q rune = 81
)

type exposeComp struct {
	cmd    *cobra.Command
	path   string
	detach bool
	untree bool
	watch  bool
}

func buildExpose() exposeComp {
	c := exposeComp{
		cmd: &cobra.Command{
			Use:           "expose",
			Short:         "Exposes all the notes in a directory",
			Args:          cobra.NoArgs,
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().BoolVarP(&c.untree, "untree", "u", false, "disable default tree file organization depending on types")
	c.cmd.Flags().BoolVarP(&c.detach, "detach", "d", false, "leaves the program without remove the files")
	c.cmd.Flags().StringVarP(&c.path, "path", "p", "", "spit it all out in the defined path")
	c.cmd.Flags().BoolVarP(&c.watch, "watch", "w", false, "start watching for changes")

	return c
}

func (c *exposeComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		var (
			box              = store.New()
			views            = box.List()
			targetDir string = config.App.Paths.CacheDir + "/"
		)

		_ = os.MkdirAll(config.App.Paths.CacheDir, os.ModePerm)

		if c.path != "" {
			if _, err := os.Stat(c.path); err != nil {
				return err
			}

			targetDir = c.path + "/" + config.AppName + "/"
		}

		if !c.untree {
			for _, g := range box.GetGroups() {
				_ = os.MkdirAll(targetDir+"/"+g, os.ModePerm)
			}
		}

		var err error

		for _, v := range views {
			var f *os.File

			if c.untree || v.Group == "" {
				f, err = os.Create(targetDir + v.Key + "-" + v.Tag)
			} else if v.Extension != "" {
				f, err = os.Create(targetDir + v.Group + "/" + v.Key + "-" + v.Tag + "." + v.Extension)
			} else {
				f, err = os.Create(targetDir + v.Group + "/" + v.Key + "-" + v.Tag)
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

		color.New(color.FgHiMagenta).Fprintln(os.Stdout, "Files exposed on "+targetDir)

		if c.detach {
			return nil
		}

		if err = keyboard.Open(); err != nil {
			return err
		}

		defer keyboard.Close()

		color.New(color.FgHiBlack).Fprintln(os.Stdout, "Click Q or Ctrl+C to exit")

		if c.watch {
			_ = filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, err error) error {
				if !d.IsDir() {
					go c.watchFile(path, box)
				}

				return err
			})
		}

		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				return err
			}

			if key == keyboard.KeyCtrlC || char == Q || char == q {
				break
			}
		}

		_ = os.RemoveAll(targetDir)

		if c.path == "" {
			_ = os.MkdirAll(config.App.Paths.CacheDir, os.ModePerm)
		}

		return nil
	}
}

func (e *exposeComp) watchFile(originalPath string, d store.SetModifier) {
	w := fswatch.NewFileWatcher(originalPath, 1)
	w.Start()

	counter := 0

	for w.IsRunning() {
		select {
		case <-w.Modified():
			if counter > 0 {
				key, _, found := strings.Cut(path.Base(originalPath), "-")
				if !found {
					w.Stop()

					break
				}

				content, err := ioutil.ReadFile(originalPath)
				if err != nil {
					w.Stop()

					break
				}

				err = d.ModifyContent(key, string(content))
				if err != nil {
					panic(err)
				}
			}
		case <-w.Moved():
			var resolvedPath string

			_ = filepath.WalkDir(config.App.Paths.CacheDir, func(p string, d fs.DirEntry, err error) error {
				if !d.IsDir() && d.Name() == path.Base(originalPath) {
					resolvedPath = p
				}

				return err
			})

			w.Stop()

			if resolvedPath == "" {
				color.New(color.FgHiYellow).Fprintln(os.Stderr, "Error: "+originalPath+" file is untrackable, unfollowed")

				return
			}

			sType := path.Base(path.Dir(resolvedPath))
			key, _, _ := strings.Cut(path.Base(resolvedPath), "-")

			if config.AppName == sType {
				sType = config.TypeMain
			}

			err := d.ModifyType(key, sType)
			if err != nil {
				if errors.Is(err, store.ErrMainAlreadyExists) || errors.Is(err, store.ErrInvalidNoteType) {
					color.New(color.FgHiYellow).Fprintf(os.Stderr, "Error: %v, the file type cannot be updated but it's still being watched\n", err)
				} else if errors.Is(err, store.ErrNoteNotFound) {
					color.New(color.FgHiYellow).Fprintf(os.Stderr, "Error: %v, it means that is untrackable, out of sight\n", err)
				} else {
					panic(err)
				}
			}

			e.watchFile(resolvedPath, d)
		}

		counter++
	}
}
