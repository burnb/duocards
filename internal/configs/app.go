package configs

import (
	"flag"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
)

const (
	Name = "duocards"

	defaultFilePath       = "text.txt"
	defaultSeparator      = "/"
	defaultNativeLanguage = "en"
)

type App struct {
	ApiToken       string `envconfig:"API_TOKEN" required:"true"`
	DeckId         string `envconfig:"DECK_ID" required:"true"`
	FilePath       string
	Separator      string
	Command        string
	NativeLanguage string
	flagSet        *flag.FlagSet
}

// Prepare variables to static configuration
func (c *App) Prepare() (err error) {
	c.flagSet = flag.NewFlagSet(Name, flag.ExitOnError)
	c.flagSet.Usage = c.printUsage

	if err = envconfig.Process("", c); err != nil {
		return err
	}

	help := c.flagSet.Bool("h", false, "Show this help message and exit")
	c.flagSet.StringVar(&c.FilePath, "i", defaultFilePath, "Input file path")
	c.flagSet.StringVar(&c.Separator, "s", defaultSeparator, "Columns separator in the input file")
	c.flagSet.StringVar(&c.NativeLanguage, "l", defaultNativeLanguage, "Your native language")
	if err = c.flagSet.Parse(os.Args[1:]); err != nil {
		return err
	}

	args := c.flagSet.Args()
	if *help || len(args) == 0 {
		c.flagSet.Usage()
		os.Exit(1)
	}
	c.Command = args[0]

	return nil
}

func (c *App) printUsage() {
	fmt.Println("\nUsage: duocards [OPTIONS] COMMAND\n\nOptions:")
	c.flagSet.PrintDefaults()
	fmt.Println(
		`
Commands:
    import               Import from file`,
	)
}
