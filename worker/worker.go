package worker

import (
	"sync"

	"humn.ai/phan/logger"
	"humn.ai/phan/mapbox"
	"humn.ai/phan/models"
)

type Pool struct {
	workersCount int
	waitGroup    *sync.WaitGroup
	mapboxClient mapbox.Client
	input        chan models.Input
	output       chan models.Output
}

func NewPool(mapboxClient mapbox.Client, inChannel chan models.Input, outChannel chan models.Output, workersCount int) *Pool {
	return &Pool{
		waitGroup:    &sync.WaitGroup{},
		workersCount: workersCount,
		mapboxClient: mapboxClient,
		input:        inChannel,
		output:       outChannel,
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.workersCount; i++ {
		p.waitGroup.Add(1)
		go runTask(p.mapboxClient, p.input, p.output, p.waitGroup)
	}
}

func (p *Pool) Wait() {
	p.waitGroup.Wait()
}

// runTask reads coordinates from an input channel, queries Mapbox API to get the corresponding postcode, then pushes
// all these info into an output channel. runTask terminates when the input channel is closed.
func runTask(mapboxClient mapbox.Client, input chan models.Input, output chan models.Output, wg *sync.WaitGroup) {
	defer wg.Done()
	for in := range input {
		postcode, err := mapboxClient.GetPostcode(in.Longitude, in.Latitude)
		if err != nil {
			logger.Errorf("failed to query postcode, pushing back to input channel: %s", err.Error())
			select {
			case input <- in:
				// Push data back to the channel if Mapbox has a temporary problem
			default:
				logger.Fatalf("failed to push to input channel: buffer overflow")
				return
			}
			continue
		}

		out := models.Output{
			Coordinate: in.Coordinate,
			Postcode:   postcode,
		}

		select {
		case output <- out:
			// Push data to output channel
		default:
			logger.Fatalf("failed to push to output channel: buffer overflow")
			return
		}
	}
	logger.Infof("Input channel is closed, nothing more to be read")
}
