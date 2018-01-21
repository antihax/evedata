# evedata


[EVEData.org website](https://www.evedata.org)

[![Build Status](https://travis-ci.org/antihax/evedata.svg?branch=master)](https://travis-ci.org/antihax/evedata)
[![codecov](https://codecov.io/gh/antihax/evedata/branch/master/graph/badge.svg)](https://codecov.io/gh/antihax/evedata)

## Contact
See @antihax on #devfleet #tweetfleet Slack.

## Contributing

You will need Docker for the mock services

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
4. Run ./test.sh

Before working on your local copy, please use a seperate branch.
If there are tests in the package, please make sure you add tests for your work, unless it will hit a public CCP service.