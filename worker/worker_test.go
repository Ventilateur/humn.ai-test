package worker

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"humn.ai/phan/models"
)

type MockedMapboxClient struct {
	mock.Mock
}

func (m *MockedMapboxClient) GetPostcode(longitude, latitude float64) (string, error) {
	args := m.Called(longitude, latitude)
	return args.String(0), args.Error(1)
}

func prepareChannels(inSize, outSize int) (chan models.Input, chan models.Output) {
	input := make(chan models.Input, inSize)
	for i := 0; i < inSize; i++ {
		input <- models.Input{
			Coordinate: models.Coordinate{
				Longitude: float64(i),
				Latitude:  float64(i),
			},
		}
	}
	output := make(chan models.Output, outSize)
	return input, output
}

func TestRunTask(t *testing.T) {
	bufSize := 4
	mockedMapboxClient := new(MockedMapboxClient)
	mockedMapboxClient.On("GetPostcode", mock.Anything, mock.Anything).Return("PC", nil)

	input, output := prepareChannels(bufSize, bufSize)

	var wg sync.WaitGroup
	wg.Add(1)
	go runTask(mockedMapboxClient, input, output, &wg)

	for i := 0; i < bufSize; i++ {
		require.Equal(t, models.Output{
			Coordinate: models.Coordinate{
				Longitude: float64(i),
				Latitude:  float64(i),
			},
			Postcode: "PC",
		}, <-output)
	}

	close(input)
}

func TestPushBack(t *testing.T) {
	bufSize := 4
	mockedMapboxClient := new(MockedMapboxClient)

	// GetPostcode will fail twice
	mockedMapboxClient.On("GetPostcode", mock.Anything, mock.Anything).Return("", errors.New("")).Once()
	mockedMapboxClient.On("GetPostcode", mock.Anything, mock.Anything).Return("", errors.New("")).Once()
	mockedMapboxClient.On("GetPostcode", mock.Anything, mock.Anything).Return("PC", nil)

	input, output := prepareChannels(bufSize, bufSize)

	var wg sync.WaitGroup
	wg.Add(1)
	go runTask(mockedMapboxClient, input, output, &wg)

	// The output should be 2, 3, 0, 1 because 0 and 1 were pushed back to the channel due to failing GetPostcode
	outValues := []int{2, 3, 0, 1}
	for i := 0; i < bufSize; i++ {
		require.Equal(t, models.Output{
			Coordinate: models.Coordinate{
				Longitude: float64(outValues[i]),
				Latitude:  float64(outValues[i]),
			},
			Postcode: "PC",
		}, <-output)
	}

	close(input)
}
