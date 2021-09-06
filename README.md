# humn.ai test

## Requirements

See [REQUIREMENTS.md](REQUIREMENTS.md)

## Build, test and run

Test with coverage report:

```shell
make test
```

Build, get some help and run:

```shell
make build
./app --help
cat coordinates.txt | ./app <api_token> -w <number_of_workers> > output.txt
```

## Implementation

This is a traditional fan-out/fan-in pattern:

```
                          |-> worker 1 ->|
stdin -> input channel -> |-> worker 2 ->|-> output channel -> stdout 
                          |-> worker 3 ->|
```

## Notes

* This program does not deal with buffered channel overflow. The default capacity for both input and output channel 
  is 10000. Any overflow will terminate the program immaturely. An elastic buffer can solve the issue.
* The program will eventually terminate if it encounters EOF input.
* The output does not contain any space, since it's compacted JSON.
