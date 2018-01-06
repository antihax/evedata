# evedata


[EVEData.org website](https://www.evedata.org)

[![Build Status](https://travis-ci.org/antihax/evedata.svg?branch=master)](https://travis-ci.org/antihax/evedata)
[![codecov](https://codecov.io/gh/antihax/evedata/branch/master/graph/badge.svg)](https://codecov.io/gh/antihax/evedata)

## Contact
See @antihax on #devfleet #tweetfleet Slack.

## Contributing

You will need:

- A MySQL, MariaDB, or Percona server.
- A redis server.
- Docker

### Services

| Service        | Description | 
| ------------- |-------------| 
| Artifice      | Task scheduler | 
| DiscordBotTemp    | Temporary hacks to provide feasibility tests for a discord bot | 
| Hammer | ESI Consumer |  
| Nail | Database store |  
| Vanguard | Web Front End|  
| ZKillboard | ZKillboard API and RedisQ Consumer |  


### Setup your environment

1. Fork this repository and clone the fork into `gopath/src/github.com/antihax`.
2. `go get -u ./...` in the repository to install dependencies.
3. Run ./mock.sh
4. For testing: run a blank redis on 127.0.0.1:6379 with no authentication ** THIS WILL GET WIPED DURING TESTING **
5. For testing: run a SQL server on 127.0.0.1:3306 with a blank root password. Import both .sql files from earlier.
6. Decompress sql.zip and install into a database called `eve`.
7. Install evedata.sql into a database called `evedata`.
8. Testing can be configured in a disposable docker env.

Before working on your local copy, please use a seperate branch.
If there are tests in the package, please make sure you add tests for your work, unless it will hit a public CCP service.