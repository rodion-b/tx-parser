# Tx Parser Service

This is a naive implementation of ETH blocks transactions parser service for a subscribed address 


## Functional requirements
-   Users need to be able to get Currently Processed ETH Block
-   Users need to be able to subscribe to recieve notifications if transactions happened with their addresses
-   Users need to be able to fetch all of their transactions
  

## Technical requirements
- **Language:** Golang  
- **Database:** Custom Inmemory Storage
- No external libraries like geth - pure JSONRPC interaction with the blockchain
- The server needs to be accessable via rest API or command line

## How to run
- `make test` runs the tests
- `make run` runs the app locally on port 8080 without docker.
- `make lint` runs the linter


## Solution notes
-   There are 3 endpoints available:
    - GET /currentblock - gets currently latest processed block
    - GET /subscribe?address=0xCFA7DA6E16be3303F09d460Ccf31D2122A2D5604 - instructs the app to start monitoring transactions for the address
    - GET .transactions?address=0xCFA7DA6E16be3303F09d460Ccf31D2122A2D5604 - retrieve all the transactions for the address in the form of json


- standard Go project layout
- clean architecture (handler->service->repository)
- decoupling of project components is achived via the usage of interfaces
- Makefile is included
- Tx parser tests to verify all works as expected
