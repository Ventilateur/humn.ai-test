package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"humn.ai/phan/models"
	"humn.ai/phan/stdio"
	"humn.ai/phan/worker"
)

type MockedMapboxClient struct {
	mock.Mock
}

func (m *MockedMapboxClient) GetPostcode(longitude, latitude float64) (string, error) {
	args := m.Called(longitude, latitude)
	return args.String(0), args.Error(1)
}

func TestRun(t *testing.T) {
	inputData := strings.Join([]string{
		`{"lat": 0.1, "lng": 0.1}`,
		`{"lat": 1.1, "lng": 1.1}`,
		`{"lat": 2.2, "lng": 2.2}`,
		`{"lat": 3.3, "lng": 3.3}`,
		`{"lat": 4.4, "lng": 4.4}`,
	}, "\n")
	outputData := strings.ReplaceAll(strings.ReplaceAll(inputData, "}", `,"postcode":"PC"}`), " ", "") + "\n"

	mockedMapboxClient := new(MockedMapboxClient)
	mockedMapboxClient.On("GetPostcode", mock.Anything, mock.Anything).Return("PC", nil)

	in := make(chan models.Input, bufferSize)
	out := make(chan models.Output, bufferSize)
	outBuffer := new(bytes.Buffer)
	writeFinished := make(chan bool)

	workerPool := worker.NewPool(mockedMapboxClient, in, out, 1)

	workerPool.Run()
	go stdio.Write(outBuffer, out, writeFinished)
	go stdio.Read(strings.NewReader(inputData), in)

	workerPool.Wait()
	close(out)
	<-writeFinished

	require.Equal(t, outputData, outBuffer.String())
}
