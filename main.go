package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/dhulihan/adafruit-io/aio"
	"os"
	// "reflect"
	"fmt"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	// a := aio.NewContext()
	// fmt.Printf("%+v\n", a)

	app := cli.NewApp()
	app.Name = "adafruit-io"
	app.Version = "1.0.0"
	app.Usage = "Send data to your adafruit.io dashboard"
	app.Flags = []cli.Flag{
		// cli.StringFlag{
		// 	Name: "verbose",
		// 	Value: "true",
		// 	Usage: "Enable to see debug messages",
		// }
		cli.StringFlag{
			Name:  "format, f",
			Value: "text",
			Usage: "Desired output format. Options: json, text (default)",
		}, cli.StringFlag{
			Name:   "key, k",
			Usage:  "Your adafruit.io secret key. $AIO_KEY tried first",
			EnvVar: "AIO_KEY",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Enable to see debug messages",
		},
	}
	app.Action = MainAction
	app.Commands = []cli.Command{
		{
			Name:        "info",
			Aliases:     []string{"i"},
			Usage:       "Get feed info",
			Description: "info <FEED ID|FEED NAME|FEED KEY>",
			Action:      InfoAction,
		},
		{
			Name:        "get",
			Aliases:     []string{"f"},
			Usage:       "Get feeds last value",
			Description: "get <FEED ID|FEED NAME|FEED KEY>",
			Action:      GetAction,
		},
		{
			Name:        "send",
			Aliases:     []string{"s"},
			Usage:       "Send a value to a feed",
			Description: "send <FEED ID|FEED NAME|FEED KEY> <VALUE>",
			Action:      SendAction,
		},
		{
			Name:    "key",
			Aliases: []string{"k"},
			Usage:   "print AIO key",
			Action:  KeyAction,
		},
	}
	Run(app)
}

func Run(app *cli.App) {
	app.Before = func(c *cli.Context) error {
		if c.GlobalString("key") == "" {
			log.Fatal("No aio key provided. Use --key KEY_HERE or export AIO_KEY=KEY_HERE")
		}

		if c.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug Mode ON")
		}
		return nil
	}
	app.Run(os.Args)
}

func MainAction(c *cli.Context) {
	log.Debug("Starting...")
	log.Debug("using adafruit.io key ", c.GlobalString("key"))
	log.Debug("Args: ", c.Args())

	feeds := aio.Feeds(c.GlobalString("key"))
	if len(feeds) > 0 {
		for _, feed := range feeds {
			fmt.Println(feed.Name)
		}

	} else {
		fmt.Println("No feeds found.")
	}

	if len(c.Args()) == 0 {
		log.Debug("No action specified")
		fmt.Println("Please provide a subcommand. Run --help for some examples.")
	}
}

func KeyAction(c *cli.Context) {
	fmt.Println(c.String("key"))
}

func InfoAction(c *cli.Context) {
	log.Debug("Args: ", c.Args())
	if len(c.Args()) == 0 {
		log.Fatal("feed id missing")
	}
	id := c.Args().First()
	for k, v := range FetchFeed(c, id) {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func FetchFeed(c *cli.Context, id string) map[string]interface{} {
	return nil
}

func GetAction(c *cli.Context) {
	log.Debug("Args: ", c.Args())
	if len(c.Args()) == 0 {
		log.Fatal("feed id missing")
	}

	a := aio.NewContext(c.GlobalString("key"))
	id := c.Args().First()
	feed, err := aio.Find(id, &a)
	if err != nil {
		log.Fatal(err)
	}
	last_value := feed.Last_Value
	if last_value != "" {
		fmt.Println(last_value)
	} else {
		log.Fatal("last_value not set")
	}
}

func SendAction(c *cli.Context) {
	log.Debug("Args: ", c.Args())
	if len(c.Args()) == 0 {
		log.Fatal("feed id missing")
	}

	if len(c.Args()) == 1 {
		log.Fatal("value is missing")
	}

	a := aio.NewContext(c.GlobalString("key"))
	id := c.Args().First()
	val := c.Args()[len(c.Args())-1]
	err := aio.Send(id, val, &a)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("OK", val)
	}
}

func FeedInfo(c *cli.Context, id string) map[string]interface{} {
	return nil
}
