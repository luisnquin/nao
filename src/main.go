package main

import "github.com/luisnquin/nao/src/cmd"

func main() {
	cmd.Execute()
}

/*
flag.String("", "", "patas usage")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	ee := viper.GetString("")

	fmt.Println(ee)
*/

// viper.AddConfigPath(appdir.New("nao").UserConfig())
