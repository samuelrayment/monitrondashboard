package monitrondashboard

// HTTP client code for the monitron dashboard

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sort"
)

type BuildUpdate struct {
	builds []build
	err    error
}

type BuildFetcher interface {
	BuildChannel() chan BuildUpdate
}

// StringUntilReader provides the ReadString method
// from bufio.Reader
type StringUntilReader interface {
	ReadString(delim byte) (line string, err error)
}

func NewBuildFetcher() BuildFetcher {
	buildFetcher := tcpBuildFetcher{}
	buildFetcher.buildChannel = make(chan BuildUpdate)
	buildFetcher.fetchBuilds()
	return buildFetcher
}

type tcpBuildFetcher struct {
	conn         net.Conn
	reader       StringUntilReader
	buildChannel chan BuildUpdate
}

func (bf tcpBuildFetcher) BuildChannel() chan BuildUpdate {
	return bf.buildChannel
}

func (bf *tcpBuildFetcher) fetchBuilds() {
	var err error
	if bf.conn, err = net.Dial("tcp", "localhost:9988"); err != nil {
		fmt.Printf("Error Connecting")
		return
	}
	bf.reader = bufio.NewReader(bf.conn)
	go bf.readLoop()
}

func (bf tcpBuildFetcher) readLoop() {
	for {
		bf.processBuilds()
	}
}

func (bf tcpBuildFetcher) processBuilds() {
	buildStatus, err := bf.reader.ReadString('\n')
	if err != nil {
		bf.buildChannel <- BuildUpdate{
			builds: []build{},
			err:    errors.New("Network Error"),
		}
	}
	var buildCollection jsonBuildCollection
	if err := json.Unmarshal([]byte(buildStatus), &buildCollection); err != nil {
		bf.buildChannel <- BuildUpdate{
			builds: []build{},
			err:    fmt.Errorf("Cannot Parse JSON: %s", err),
		}
		return
	}
	builds := bf.processJSONBuildIntoBuildList(buildCollection)
	bf.buildChannel <- BuildUpdate{
		builds: builds,
		err:    nil,
	}
}

// sortByName is a sort interface for a []build that sorts by build name
type sortByName []build

func (s sortByName) Len() int {
	return len(s)
}

func (s sortByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortByName) Less(i, j int) bool {
	return len(s[i].name) < len(s[j].name)
}

func (bf tcpBuildFetcher) processJSONBuildIntoBuildList(buildCollection jsonBuildCollection) []build {
	addBuildsFromSet := func(buildSet []jsonBuild, state buildState, buildList []build) []build {
		for _, i := range buildSet {
			buildList = append(buildList,
				build{
					name:         i.Name,
					buildState:   state,
					building:     i.Building,
					acknowledger: i.Acknowledger,
				})
		}
		return buildList
	}

	returnBuilds := addBuildsFromSet(buildCollection.Failing,
		BuildStateFailed, []build{})
	returnBuilds = addBuildsFromSet(buildCollection.Acknowledged,
		BuildStateAcknowledged, returnBuilds)
	returnBuilds = addBuildsFromSet(buildCollection.Healthy,
		BuildStatePassed, returnBuilds)

	sort.Sort(sortByName(returnBuilds))
	return returnBuilds
}

type jsonBuildCollection struct {
	Failing      []jsonBuild `json:"failing"`
	Acknowledged []jsonBuild `json:"acknowledged"`
	Healthy      []jsonBuild `json:"healthy"`
}

type jsonBuild struct {
	Name         string `json:"name"`
	Building     bool   `json:"building"`
	Acknowledger string `json:"user"`
}
