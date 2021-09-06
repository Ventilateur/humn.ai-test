package stdio

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"humn.ai/phan/models"
)

var dataRead = []string{
	"{\"lat\": 0.0, \"lng\": 0.0}",
	"{\"lat\": 1.0, \"lng\": 1.0}",
	"{\"lat\": 2.0, \"lng\": 2.0}",
	"{\"lat\": 3.0, \"lng\": 3.0}",
	"{\"lat\": 4.0, \"lng\": 4.0}",
}

var dataWrite = []models.Output{
	{
		models.Coordinate{
			Longitude: 0.0,
			Latitude:  0.0,
		},
		"A",
	},
	{
		models.Coordinate{
			Longitude: 1.0,
			Latitude:  1.0,
		},
		"B",
	},
	{
		models.Coordinate{
			Longitude: 2.0,
			Latitude:  2.0,
		},
		"C",
	},
}

func TestReader(t *testing.T) {
	r := strings.NewReader(strings.Join(dataRead, "\n"))
	c := make(chan models.Input, len(dataRead))
	Read(r, c)

	for i := 0; i < len(dataRead); i++ {
		actual := <-c
		require.Equal(t, float64(i), actual.Longitude)
		require.Equal(t, float64(i), actual.Latitude)
	}
}

func TestWriter(t *testing.T) {
	buf := make(chan models.Output, len(dataWrite))
	for _, data := range dataWrite {
		buf <- data
	}
	finished := make(chan bool)

	out := new(bytes.Buffer)
	go Write(out, buf, finished)

	time.Sleep(time.Second)
	close(buf)
	<-finished

	expected := strings.Join([]string{
		"{\"lat\":0,\"lng\":0,\"postcode\":\"A\"}",
		"{\"lat\":1,\"lng\":1,\"postcode\":\"B\"}",
		"{\"lat\":2,\"lng\":2,\"postcode\":\"C\"}",
	}, "\n") + "\n"

	require.Equal(t, expected, out.String())
}
