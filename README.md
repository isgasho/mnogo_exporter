# mnogo exporter

This is the new MongoDB exporter implementation that handles ALL metrics exposed by MongoDB monitoring commands.
Currently, these metric sources are implemented:
- $collStats
- getDiagnosticData
- replSetGetStatus
- serverStatus

## Testing 
### Initialize tools and dependencies
In order to install tools to format, test and build the exporter, you need to run this command:
```
make init
```
It will install gomports, goreleaser, golangci-lint and reviewdog.

### Starting the sandbox
The testing sandbox starts n MongoDB instances as follow:
- 3 Instances for shard 1 at ports 17001, 17002, 17003
- 3 instances for shard 2 at ports 17004, 17005, 17006
- 3 config servers at ports 17007, 17008, 17009
- 1 mongos server at port 17000
- 1 stand alone instance at port 27017
All instances are currently running without user and password so for example, to connect to the **mongos** you can just use:
```
mongo mongodb://127.0.0.1:17001/admin
```
The sandbox can be started using the provided Makefile using: `make test-cluster` and it can be stopped using `make test-cluster-clean`

To run the tests, just run `make test`



