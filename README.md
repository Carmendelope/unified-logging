# Unified logging

`unified-logging` implements a distributed system to collect logs from user applications running on the Nalej platform. It uses standard components for log ingestion, storage and querying, adds an aggregation layer and leverages the existing platform constructs for the distributed nature.


## Getting Started

## Architecture

Kubernetes, and therefore the Nalej platform, uses the standard Docker mechanism to collect logs from applications: everything a container writes to `stdout` or `stderr` gets captured in a file on the node where that container runs. The Unified Logging solution deploys a filebeat instance on each node to scrape those logs, filter them and store them in a cluster-local ElasticSearch instance.

First of all, filebeat annotates each log line with Kubernetes information. Then, it discards any log line that does not originate from a container with the label `nalej-organization` set - in other words, it only deals with user applications. Next, it discards lines originating from `zt-sidecar` containers, as though they are deployed in a user namespace, they are not part of the user logging infromation. Lastly, it drops almost all annotations for a log line except for the Kubernetes namespace and labels, to save space.

Currently the standard filebeat/ElasticSearch indexing is used - this should be changed to an index per namespace, as this will be much more efficient; this is what we will query on (as this is per application/fragment).

An application cluster-local component, `unified-logging-slave`, implements `Search` and `Expire` endpoints to retrieve cluster-local logs. Search implements filters for application instance and service group instance as well as free text search and an optional time range. Expire will delete all logs for a specific application instance.

The logging slave will return at most 10,000 log lines (this is a default ElasticSearch limitation). To retrieve more logs, we should implement the use of the Elastic scroll or pagination APIs.

On the management cluster, the `unified-logging-coord` implements the same `Search` and `Expire` endpoints, except that it executes them on the relevant application clusters (currently, just on all available). When all logs are retrieved, the coordinator merges and sorts them before returning.

The end-to-end mechanism follows our standard architecture of Public API -> Coordinator -> Application cluster API -> Slave.

## To do

As per above:

- Optimization of cluster-local queries by reorganizing the storage indexing
- Optimization of retrieval by only querying relevant clusters and parallel querying.
- Pagination or scroll API use
- Expiration for time range instead of all logs for an instance
- Potentially storing certain log lines (by filter? with errors or warnings?) on the management cluster for longer term storage / disaster recovery and analysis.

### Prerequisites

#### Slave

`unified-logging-slave` depends on ElasticSearch running locally, without any security mechanism. Furthermore, it expects filebeat to ingest the logs in ElasticSearch. To this end, we have deployments for both as part of the unified logging package.

```
$ ./unified-logging-slave run --help
Launch the server API

Usage:
  unified-logging-slave run [flags]

Flags:
      --elasticAddress string   ElasticSearch address (host:port) (default "localhost:9200")
  -h, --help                    help for run
      --port int                Port for Unified Logging Slave gRPC API (default 8322)

Global Flags:
      --consoleLogging   Pretty print logging
      --debug            Set debug level
```

#### Coordinator

```
$ ./unified-logging-coord run --help
Launch the server API

Usage:
  unified-logging-coord run [flags]

Flags:
      --appClusterPort int          Port used by app-cluster-api (default 443)
      --appClusterPrefix string     Prefix for application cluster hostnames (default "appcluster")
      --caCert string               Alternative certificate file to use for validation
  -h, --help                        Help for run
      --skipServerCertValidation    Don't validate TLS certificates
      --port int                    Port for Unified Logging Coordinator gRPC API (default 8323)
      --systemModelAddress string   System Model address (host:port) (default "localhost:8800")
      --useTLS                      Use TLS to connect to application cluster (default true)

Global Flags:
      --consoleLogging   Pretty print logging
      --debug            Set debug level
```

### Build and compile

In order to build and compile this repository use the provided Makefile:

```
make all
```

This operation generates the binaries for this repo, download dependencies,
run existing tests and generate ready-to-deploy Kubernetes files.

### Run tests

Tests are executed using Ginkgo. To run all the available tests:

```
make test
```

### Integration tests

The following table contains the variables that activate the integration tests

 | Variable  | Example Value | Description |
 | ------------- | ------------- |------------- |
 | RUN_INTEGRATION_TEST  | true | Run integration tests |
 | IT_ELASTIC_ADDRESS  | localhost:9200 | ElasticSearch Address |

To run Elastic: `docker run --rm -it -p 9200:9200 docker.elastic.co/elasticsearch/elasticsearch-oss:6.6.0 elasticsearch`


### Update dependencies

Dependencies are managed using Godep. For an automatic dependencies download use:

```
make dep
```

In order to have all dependencies up-to-date run:

```
dep ensure -update -v
```

## User client interface

### API

All endpoints implement:
- `Search` with a `SearchRequest` as argument and a `LogResponse` as response, and
- `Expire` with an `ExpirationRequest` as argument and a `common.Success` (true or false) as response.

Common for both requests are an organization ID and an application instance ID. On top, a `SearchRequest` also has fields for a service group ID, a log message free text filter string, a time range and a sort order.

The `LogResponse` returns the organization ID and application instance ID, the actual time range of the log lines returned and an array of timestamp / message tuples.

See [unified-logging](https://github.com/nalej/grpc-protos/tree/master/unified-logging) for details.

### CLI

The public API CLI only implements the search request, as follows:

```
$ ./public-api-cli log search --help
Search application logs based on application and service group instance

Usage:
  public-api-cli log search [filter string] [flags]

Flags:
      --asc                   Sort results in ascending time order
      --desc                  Sort results in descending time order
      --from string           Start time of logs
  -h, --help                  help for search
      --instanceID string     Application instance identifier
      --redirectResultAsLog   Redirect the result to the CLI log
      --sgInstanceID string   Service group instance identifier
      --to string             End time of logs

Global Flags:
      --cacert string             Path of the CA certificate to validate the server connection
      --consoleLogging            Pretty print logging
      --debug                     Set debug level
      --skipServerCertValidation  Use a insecure connection to connect to the server
      --nalejAddress string       Address (host) of the Nalej platform
      --organizationID string     Organization identifier
```

## Contributing

Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.


## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/nalej/unified-logging/tags). 

## Authors

See also the list of [contributors](https://github.com/nalej/unified-logging/contributors) who participated in this project.

## License
This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.





## Usage




