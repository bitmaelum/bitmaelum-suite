# Build

We are using `make` to build our suite.

The suite itself consists of several main apps and some additional tools. The tools are not needed for running BitMaelum, but can come in handy during debugging and development.

To test the current suite:

    make test
    
To build the current suite on your own platform and archicture:

    make build
    
The binaries will be compiled into `release` directory.

If you want to build a specific app or tool, just supply the name:

    make bm-server
    
This will rebuild only the `bm-server` binary.
 

## Cross platform builds
It's possible to build different platforms and architectures:

    make build-all
    
All binaries will be compiled into the `release/<os>-<arch>` directory.
    
Currently, the following platforms are automatically build:

   - windows-amd64 
   - linux-amd64 
   - darwin-amd64

It's possible to compile only a single platform:

    make linux-amd64

It's not possible to build cross-compile a single binary (ie: `make linux-amd bm-server` will not work) 
