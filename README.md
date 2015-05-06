Monitron Dashboard
==================

A Dashboard for the Monitron build server, using only the terminal for presentation; written using 
the excellent [Termbox](https://github.com/nsf/termbox-go).

[![Build Status](https://travis-ci.org/samuelrayment/monitrondashboard.svg?branch=master)](https://travis-ci.org/samuelrayment/monitrondashboard)

Getting Started
---------------

We assume you have a working go environment, obtainable here: http://golang.org/.

* go get github.com/samuelrayment/monitrondashboard
* go install github.com/samuelrayment/monitrondashboard

You can now run the dashboard using:

    monidash -a <hostname:port for your monitron server>

You can also provide the address using an environment variable: 

    MD_ADDRESS

Docker
------

You can run the monitron dashboard from a docker image, you must use a tty and interactive mode
though adjusting the MD_ADDRESS to the correct monitron server address, e.g.:

    sudo docker run -ti --rm --name monidash -e MD_ADDRESS=192.168.0.13:9988 bestriped/monitrondashboard
