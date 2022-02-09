# Discobeat

Discobeat is an elastic [beat](https://www.elastic.co/products/beats) that publishes messages from Discord to 
elasticsearch

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/fergloragain/discobeat`

## Getting Started with Discobeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Build

To build the binary for Discobeat run the command below. This will generate a binary
in the same directory with the name discobeat.

```
make
```

## Configuration

Discobeat can be used with either a standard Discord username and password, or with a Bot token. 

For each user specified, Discord will query all guilds and all channels that user has access to, and publish the most  recent messages in each channel to elasticsearch. 

If `archive` is set to `true`, Discobeat will fetch the entire message history for each channel, 100 messages at a time,  though this will result in a much longer execution time.  

If a list of Guild IDs is specified under a particular user, then only those Guilds will be queried for messages. If no Guilds are specified, all Guilds associated with the user will be queried. Furthermore, if a list of channel IDs is specified under a Guild, only those channels within that Guild will be queried. If a list of channel IDs is not specified under a Guild, all channels under the Guild will be queried. 

For an example of the configuration structure, see [this reference configuration](/discobeat.reference.yml)
 
### Run

To run Discobeat with debugging output enabled, run:

```
./discobeat -c discobeat.yml -e -d "*"
```

### Test

To test Discobeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```

### Cleanup

To clean  Discobeat source code, run the following command:

```
make fmt
```

To clean up the build directory and generated artifacts, run:

```
make clean
```

### Clone

To clone Discobeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/fergloragain/discobeat
git clone https://github.com/fergloragain/discobeat ${GOPATH}/src/github.com/fergloragain/discobeat
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires 
[docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following 
command:

```
make release
```

This will fetch and create all images required for the build process. The whole process to finish can take several 
minutes.
