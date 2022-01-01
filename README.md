# evedata


[EVEData.org website](https://www.evedata.org)


[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/O5O33VK5S)

## Contact

See @antihax on #devfleet #tweetfleet Slack.

## Contributing

You will need Docker for the mock services

### Services

| Service        | Description | 
| ------------- |-------------| 
| Artifice      | Task scheduler | 
| Conservator    | Integration (discord, slack, ts3, mumble) | 
| Hammer | Main ESI Consumer | 
| KillmailDump | Dumps killmail stream to json files |   
| Nail | Database store |  
| Squirrel | Not used yet. Pulls static data into DB get updates faster. |  
| TokenServer | CCP OAuth2 Caching service | 
| Vanguard | Web Front End|  
| ZKillboard | ZKillboard API and RedisQ Consumer |  


### Setup your environment

1. Fork this repository and clone the fork into `gopath/src/github.com/antihax`.
2. `go get -u ./...` in the repository to install dependencies.
3. Run ./mock.sh
4. Run ./test.sh

Before working on your local copy, please use a seperate branch.
If there are tests in the package, please make sure you add tests for your work, unless it will hit a public CCP service.
