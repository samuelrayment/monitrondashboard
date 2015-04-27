package monitrondashboard

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var testData string = "{\"type\":\"builds\",\"error\":\"\",\"failing\":[{\"name\":\"Failing Build\",\"building\":false,\"user\":\"\",\"url\":\"http://localhost:8000/job/Failing%20Build/\",\"number_of_failures\":1,\"failing_since\":1425590828000}],\"acknowledged\":[],\"healthy\":[{\"name\":\"Build\",\"building\":false,\"user\":\"\",\"url\":\"http://localhost:8000/job/Test/\",\"number_of_failures\":0,\"failing_since\":0}]}"

type MockStringUntilReader struct {
	mock.Mock
}

func (msur *MockStringUntilReader) ReadString(delim byte) (line string, err error) {
	args := msur.Called(delim)
	return args.String(0), args.Error(1)
}

func TestProcessBuildsCreatesAndSortsABuildSlice(t *testing.T) {
	mockStringReader := MockStringUntilReader{}
	buildFetcher := tcpBuildFetcher{
		conn:         nil,
		reader:       &mockStringReader,
		buildChannel: make(chan BuildUpdate, 2),
	}
	mockStringReader.Mock.On("ReadString", '\n').Return(testData, nil)

	buildFetcher.processBuilds()

	buildUpdate := <-buildFetcher.buildChannel
	if buildUpdate.err != nil {
		t.Fatalf("processBuilds() should not send an error when given good data.")
	}

	if len(buildUpdate.builds) != 2 {
		t.Fatalf("Expected processBuilds() => %d builds, expected 2.", len(buildUpdate.builds))
	}

	// The builds should be sorted alphabetically
	firstBuild := buildUpdate.builds[0]
	assert.Equal(t, "Build", firstBuild.name, "Builds should be sorted alphabetically so 'Build' is first")
	assert.Equal(t, false, firstBuild.building)
	assert.Equal(t, "", firstBuild.acknowledger)
	assert.Equal(t, BuildStatePassed, firstBuild.buildState)

	secondBuild := buildUpdate.builds[1]
	assert.Equal(t, false, secondBuild.building)
	assert.Equal(t, "", secondBuild.acknowledger)
	assert.Equal(t, BuildStateFailed, secondBuild.buildState)
	assert.Equal(t, "Failing Build", secondBuild.name, "Builds should be sorted alphabetically so 'Failing Build' is second")
}

func TestProcessBuildsErrorsIfItCantParseJSON(t *testing.T) {
	mockStringReader := MockStringUntilReader{}
	buildFetcher := tcpBuildFetcher{
		conn:         nil,
		reader:       &mockStringReader,
		buildChannel: make(chan BuildUpdate, 2),
	}
	// Return bad data
	mockStringReader.Mock.On("ReadString", '\n').Return("{\"a\"}", nil)

	buildFetcher.processBuilds()

	buildUpdate := <-buildFetcher.buildChannel
	assert.Error(t, buildUpdate.err, "processBuilds() should error on malformed json")
}

func TestProcessBuildsErrorsOnNetworkError(t *testing.T) {
	mockStringReader := MockStringUntilReader{}
	buildFetcher := tcpBuildFetcher{
		conn:         nil,
		reader:       &mockStringReader,
		buildChannel: make(chan BuildUpdate, 2),
	}
	// Return bad data
	mockStringReader.Mock.On("ReadString", '\n').Return("", errors.New("network error"))

	buildFetcher.processBuilds()

	buildUpdate := <-buildFetcher.buildChannel
	assert.Error(t, buildUpdate.err, "processBuilds() should error on a network error")
}
