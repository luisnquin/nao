package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/goccy/go-json"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type ConfigCmd struct {
	*cobra.Command

	config *config.Core
	log    *zerolog.Logger
	list   bool
}

func BuildConfig(log *zerolog.Logger, config *config.Core) ConfigCmd {
	c := ConfigCmd{
		Command: &cobra.Command{
			Use:           "config",
			Short:         "Allows you to move your settings from the cli",
			Args:          cobra.MaximumNArgs(3),
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		config: config,
		log:    log,
	}

	c.RunE = LifeTimeDecorator(log, "config", c.Main())

	c.Flags().BoolVarP(&c.list, "list", "l", false, " List all variables set in config file, along with their values")

	log.Trace().Msg("the 'config' command has been created")

	return c
}

func (c *ConfigCmd) Main() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		const maxNbOfArgs = 2

		if c.list {
			result, err := yaml.Marshal(c.config)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "%s\n", result)

			return nil
		}

		if l := len(args); l == 0 {
			return config.InitPanel(c.config)
		} else if l > maxNbOfArgs {
			return fmt.Errorf("accepts at most %d arg(s), received %d", maxNbOfArgs, len(args))
		} else if l == 1 {
			value, err := NavigateMapAndGet(c.getConfigAsMap(), args[0])
			if err != nil {
				return err
			}

			fmt.Fprintln(os.Stdout, value)

			return nil
		}

		configMap := c.getConfigAsMap()

		err := NavigateMapAndSet(configMap, args[0], args[1]) // TODO: validation
		if err != nil {
			return err
		}

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(configMap)

		if err := json.NewDecoder(b).Decode(c.config); err != nil {
			return err
		}

		return c.config.Save()
	}
}

func (c ConfigCmd) getConfigAsMap() map[string]any {
	b := new(bytes.Buffer)

	if err := json.NewEncoder(b).Encode(c.config); err != nil {
		panic(err)
	}

	var result map[string]any

	if err := json.NewDecoder(b).Decode(&result); err != nil {
		panic(err)
	}

	return result
}
