package main

import (
	"log"
	"os"

	"notionsync/tools/notion"
	"notionsync/tools/todo"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "sync",
		Usage: "todo sync notion!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "notionSecret",
				Aliases:  []string{"ns"},
				Usage:    "notion secret",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "notionDatabaseID",
				Aliases:  []string{"nd"},
				Usage:    "notion databaseID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "todoClientID",
				Aliases:  []string{"tc"},
				Usage:    "todo clientID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "todoClientSecret",
				Aliases:  []string{"tcs"},
				Usage:    "todo client secret",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			var (
				notionSecret     = c.String("notionSecret")
				notionDatabaseID = c.String("notionDatabaseID")
				todoClientID     = c.String("todoClientID")
				todoClientSecret = c.String("todoClientSecret")
			)

			notionAPI := notion.New(notionSecret, notionDatabaseID)
			todoAPI, err := todo.New(todoClientID, todoClientSecret, notionAPI)
			if err != nil {
				panic(err)
			}

			if err := todoAPI.UpdateNotionAllToDo(); err != nil {
				panic(err)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
