package main

import (
	"github.com/codegangsta/cli"
	md "github.com/samuelrayment/monitrondashboard"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "monidash"
	app.Usage = "Terminal based dashboard for the Monitron 5000"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "address, a",
			Usage:  "Address for the Monitron server to connect to.",
			EnvVar: "MD_ADDRESS",
		},
	}
	app.Action = mainAppAction
	app.Run(os.Args)
}

func mainAppAction(c *cli.Context) {
	if c.String("address") == "" {
		log.Printf("You must provide the address of a server to connect to.")
		return
	}

	fetcher := md.NewBuildFetcher(c.String("address"))
	dashboard := md.NewDashboard(fetcher, md.TermboxCellDrawer{})
	dashboard.Run()
}
