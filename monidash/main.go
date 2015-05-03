package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	md "github.com/samuelrayment/monitrondashboard"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "monidash"
	app.Usage = "Terminal based dashboard for the Monitron 5000"
	app.Version = "0.1.0"
	app.Action = func(c *cli.Context) {
		fmt.Printf("Monitron 5000\n")
		dashboard := md.NewDashboard(md.TermboxCellDrawer{})
		dashboard.Run()
	}

	app.Run(os.Args)
}
