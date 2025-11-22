package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	gendiff "code"

	"github.com/urfave/cli/v3"
)

func main() {

	//var test []interface{}
	//test = []interface{}{123, []string{"vase", "rrr"}}
	//fmt.Println(test[])

	cmd := &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		UsageText: "gendiff [options] <file1> <file2>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "houtput format (default: 'stylish')",
				Value:   "stylish",
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {

			args := cmd.Args().Slice()

			if len(args) < 2 {
				return errors.New("error: requires exactly 2 file paths\nExample: gendiff file1.json file2.json \nIf you want to see help : gendiff -h")
				//return nil
			}

			format := cmd.String("format")
			file1 := args[0]
			file2 := args[1]

			diff, err := gendiff.GendDiff(file1, file2, format)

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(diff)

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
