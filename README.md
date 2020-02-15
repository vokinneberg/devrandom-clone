# /dev/random clone

Features:

* Emulates `$ cat /dev/random` shell command by eternally producing random bites sequence using [ANU Quantum Random Numbers Server](http://qrng.anu.edu.au) as the source of entropy. The random numbers are generated in real-time in their lab by measuring the quantum fluctuations of the vacuum.

Package structure:

* `entropy_source` - contains implementation and test of entropy source abstraction
* `main.go` - application entry point
* `main_test.go` - application entry point tests 
* `go.mod` - go modules configuration
* `Makefile` - make utility configuration
* `README.md` - this file


## Configuration

Application configured with constant hardcoded values. It is easy to extend to environment variables though.

List of QRNG settings:

```go
	EntropySourceUrl   = "https://qrng.anu.edu.au/API/jsonI.php" // QRNG server URL
	EntropyDataType    = "hex16" // Data format. unit8, unit16 also available.
	EntropyArrayLength = "1024" // Array length up to 1024.
	EntropyBlockSize   = "1024" // Blocks count up to 1024.
```

## Getting Up and Running

### Package manager

Go modules `go mod` used to set up dependencies for the project. Sources provided with vendor folder, so no need to download packages. Sources might be compiled right after download.

### Linters

I use [golanci-lint](https://golangci.com/ ) to statically check source code. If you're on Mac `golangci-lint` utility might be easily installed with with `brew`:

```shell
brew install golangci/tap/golangci-lint
```

For another platforms instructions could be found on the [linter official page]( https://golangci.com/).

Linter could be run against all sources in the folder recursively with the following command:

```shell
make lint
```

Considering that linting tool could not be installed on target computer this option is not part of building process.

### Test

Run following command to run application unit tests:

```shell
make test
```

### Build

Run following command to build application binary file:

```shell
make build
```

Or full build process including clean up and unit tests:

```shell
make
```

### Install

Compiled binary could be installed in to `$GOBIN` folder with the following command:

```shell
make install
```

### Run

After installation library classes could be references from any ruby code with command:

```shell
make run
```

## Technical details

I decided to implement `/dev/random` using [ANU Quantum Random Numbers Server](http://qrng.anu.edu.au) as the source of entropy as it looks elegant to get entropy by measuring the quantum fluctuations of the vacuum. 
Though it would not work in absence of internet connection and not very reliable because ANU Quantum Random Numbers Server seems to be a reliable source of entropy as it's continuously tested with various randomness test algorithms including [Diehard](http://qrng.anu.edu.au/Diehard.php#). 

I implemented the entropy source as an abstraction: interface `EntropySource` that potentially could be implemented using any entropy source. `EntropySource` interface provides only one method ```Entropy(entropy chan []byte, err chan error)```
which suppose to produce random numbers eternally using entropy channel, or produce error using err channel in case of any kind of internal error or lack of entropy. Thus you can use `devrandom-clone` package in your solution either using QRNG entropy
source implementation or implementing your own entropy source.
