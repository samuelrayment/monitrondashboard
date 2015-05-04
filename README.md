Monitron Dashboard
==================

A Dashboard for the Monitron build server, using only the terminal for presentation; written using 
the excellent [Termbox](https://github.com/nsf/termbox-go).

Getting Started
---------------

We assume you have a working go environment, obtainable here: http://golang.org/.

* go get github.com/samuelrayment/monitrondashboard
* go install github.com/samuelrayment/monitrondashboard

You can now run the dashboard using:

    monidash -a <hostname:port for your monitron server>

You can also provide the address using an environment variable: 

    MD_ADDRESS

