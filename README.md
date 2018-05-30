# pflag Source for [go-config](https://github.com/micro/go-config)

The pflag source incorporates a `pflag.FlagSet` into the config tree. Why? Because I find go-config useful in non-go-micro scenarios, including cli tools created with [cobra](https://github.com/spf13/cobra#working-with-flags). This repo was born to make that experience work with the [spf13/pflag](https://github.com/spf13/pflag) package.

## Format

We expect the use of the `spf13/pflag` package. Upper case flags will be lower cased. Delimiter is determined by the `WithDelimiter` option. By default, dashes will be used as delimiters.

### Example

```
dbAddress := pflag.String("database_address", "127.0.0.1", "the db address")
dbPort := pflag.Int("database_port", 3306, "the db port)
```

Becomes

```json
{
    "database": {
        "address": "127.0.0.1",
        "port": 3306
    }
}
```

## New Source & Example
Since a `pflag.FlagSet` is required, the recommended approach is to initialize one within `PreRun` or `Run` lifecycle hooks of a command.

```go
var conf config.Config

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		conf = config.NewConfig()
		conf.Load(flagsrc.NewSource(cmd.Flags()))
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(conf.Get("consul", "addr").String("not present"))
		log.Println(conf.Get("root", "addr").String("not present"))
	},
}

func init() {
	serveCmd.Flags().StringP("consul-addr", "a", "localhost:8500", "consul address")
	rootCmd.PersistentFlags().StringP("root-addr", "r", "example.com", "the rootcmd string")
	rootCmd.AddCommand(serveCmd)
}
```
