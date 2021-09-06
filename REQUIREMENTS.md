# Humn Cloud Engineering Practical Exercise

Thank you for undertaking our practical exercise we appreciate the time you are willing to put in and look forward to seeing what you have produced. Here at Humn a lot of what we do involves processing streams of data from various sources and the exercise has been designed to simulate some of this in a slightly simplified manner. Please if you have any questions don't hesitate to ask!

You should spend no more than 4 hours on the exercise and we would like you to focus on the showing us following:
- Your approach to solving the problem. 
- Good code readability.
- A simple solution.

## Prerequisites

You will ned a computer with the following installed to undertake this task
- Go
- Bash shell or similar

You will also need to create a free MapBox account using https://account.mapbox.com/auth/signup/. Once created they will give you a mapbox API token to use for accessing the API.

## Aim

We would like you to create a fully working command line program written in Go that demonstrates a worker pool design pattern using goroutines and channels. If you have not implemented one before a worker pool is a group of goroutines running concurrently that perform the same operation with the aim of speeding up the operations being done.

The program should read a stream of JSON encoded coordinates (latitude & longitude) from standard input (stdin), process the messages using the worker pool before writing each message to standard output (stdout). 

The processing to be performed by each worker in the pool is to do a reverse geocode lookup on the coordinates from each input message using the mapbox geocoding HTTP REST API https://docs.mapbox.com/api/search/#geocoding getting the corresponding postcode for each coordinate pair. See the below Mapbox section for further guidance on using the API.

The file `coordinates.txt` has been provided containing the messages we wish you to process.

The output messages should be in the below format and output 1 per line.

```json
{ "lat": <float64>, "lng": <float64>, "postcode": <string> }
```

You should implement your service using the following design criteria:
- All the messages in the input stream must be processed by the program and output to stdout.
- The output order does not have to be the same as the input order.
- Stdin should be read by a separate single goroutine and not directly by the worker pool.
- Stdout should be written to from a separate single goroutine and not directly from the worker pool.
- You should not use a third party library to implement your worker pool.
- The mapbox API token should be provided via a mandatory command line argument.
- The number of goroutines should be configurable via a command line flag. The default should be 5. 
- Any logging performed by the program should be directed to standard error (stderr)
- The postcode returned from mapbox should be sanitised to ensure to conforms to the expected postcode format.


The following command should be used to ensure you program is working as expected

```
 cat coordinates.txt | ./your-program "api token" "pool size flag" > output.txt
```

## Deliverables

Once you have finished we would like you to return us the following:

- Application source code and tests
- A txt file containing the processed output.
- A README.md describing
  - How to run any tests you may have written.
  - Any design decisions, compromises or assumptions you have made to complete the task. 

These can either be delivered via email or via a Git repo hosted on your favourite git hosting website.


## Mapbox

You should perform the reverse geocode lookup using the `mapbox.places` part of the API. Perform the look up using the postcode type limiting the returned results to 1.

```
curl "https://api.mapbox.com/geocoding/v5/mapbox.places/<long,lat>.json?types=postcode&limit=1&access_token=YOUR_MAPBOX_ACCESS_TOKEN"
```

From the response https://docs.mapbox.com/api/search/geocoding/#geocoding-response-object the relevant field you should obtain is the `text` field from the single returned Feature.
