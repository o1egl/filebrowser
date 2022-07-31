package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/filebrowser/filebrowser/log"
)

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is filebrowser.yaml)")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug|info|warn|error|fatal)")
	rootCmd.PersistentFlags().String("log-format", "plain", "log format (plain|json)")
	rootCmd.PersistentFlags().String("log-out", "stdout", "log out (stdout|file)")
	rootCmd.PersistentFlags().String("log-file-name", "./var/filebrowser.log", "the file to write logs to")
	rootCmd.PersistentFlags().Int("log-file-age", 1, "maximum number of days to retain old log files")
	rootCmd.PersistentFlags().Int("log-file-backups", 5, "the maximum number of old log files to retain")
	rootCmd.PersistentFlags().Bool("log-file-compress", false, "determines if the rotated log files should be compressed")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "filebrowser",
	Short: "Web based file manager",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/filebrowser/")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("filebrowser")
	}

	viper.SetEnvPrefix("")
	r := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(r)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	bindFlags(rootCmd, "")
}

func bindFlags(cmd *cobra.Command, prefix string) {
	var flagPrefix string
	if cmd.HasParent() {
		flagPrefix = cmd.Name()
		if prefix != "" {
			flagPrefix = prefix + "." + cmd.Name()
		}
	}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		key := flagPrefix
		if flagPrefix != "" {
			key += "."
		}
		key += f.Name
		r := strings.NewReplacer("-", ".")
		key = r.Replace(key)
		if err := viper.BindPFlag(key, f); err != nil {
			panic(err)
		}
	})
	for _, subCmd := range cmd.Commands() {
		bindFlags(subCmd, prefix)
	}
}

func NewLogger() (log.Logger, error) {
	level, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		return nil, err
	}
	format, err := log.ParseFormat(viper.GetString("log.format"))
	if err != nil {
		return nil, err
	}
	out, err := loggerOut()
	if err != nil {
		return nil, err
	}
	logger, err := log.NewLogger(log.Configuration{
		LogLevel: level,
		Format:   format,
		Output:   out,
	})
	if err != nil {
		return nil, err
	}
	log.DefaultLogger = logger
	return logger, nil
}

func loggerOut() (log.WriteSyncer, error) {
	out := viper.GetString("log.out")
	switch out {
	case "stdout":
		return os.Stdout, nil
	case "file":
		fileWriter := log.NewFileWriter(log.FileWriterConfig{
			Filename:   viper.GetString("log.file.name"),
			MaxSize:    viper.GetInt("log.file.size"),
			MaxAge:     viper.GetInt("log.file.age"),
			MaxBackups: viper.GetInt("log.file.backups"),
			Compress:   viper.GetBool("log.file.compress"),
		})
		return fileWriter, nil
	default:
		return nil, fmt.Errorf("unsupported log out %s", out)
	}
}
