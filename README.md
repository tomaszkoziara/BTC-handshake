# BTC-handshake

An implementation of the Bitcoin handshake.

# Prerequisites

* Git clone btcd: https://github.com/btcsuite/btcd.git
* Build btcd image by running the following command in the btcd repository.
```
docker build . -t local/btcd
```

# Running the node

Run `docker compose up` to set up a dependency node. This will create a test node running on Regnet network and listening on port 18444.

Then run `go run ./src/cmd/main.go` to execute the developed node that will connect to the dependency node and perform a handshake.

# TODO

* Tests, tests, tests.
* Thorough checks on version message.
* Fill version message with all the relevant information.
* Is it over engineered?
* How does the performance change in going more barebone?
* Discovery of the nodes, or configuration to connect to multiple nodes.
* Handle error messages (e.g. if the node is non-adherent to the protocol).