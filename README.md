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
- A [mock ESI server](https://github.com/antihax/mock-esi) running on localhost:8080.
- Prometheus Server [Optional]
- Grafana Server [Optional]

### Setup your environment

1. Fork this repository and clone the fork into `gopath/src/github.com/antihax`.
2. `go get -v` in the repository to install dependencies.
3. Decompress sql.zip and install into a database called `eve`.
4. Install evedata.sql into a database called `evedata`.
5. Create a user with access to both databases.
6. Copy `config/config-example.conf` to `config/config.conf`
7. Complete the configuration.
8. To complete the OAuth2 you need three configurations from CCP.
    1. Visit https://developers.eveonline.com/
    2. Create an SSO only application, make sure the returnURL is your public address with the same path as the config. Put key and secret into the config.
    3. Create a Token application with all scopes. Make sure the return URL is similar as above and also put key and secret in config.
    4. Create a bootstrap application with all scopes. Enter into config.
    5. `go run eve-dataserver.go` and visit `http://yourpublicip:3000/X/boostrapEveAuth`.
    6. Log in with a charater.
    7. Copy the token information into the config.
9. Run the mock-esi server on 127.0.0.1:8080
10. For testing: run a blank redis on 127.0.0.1:6379 with no authentication ** THIS WILL GET WIPED DURING TESTING **
11. For testing: run a SQL server on 127.0.0.1:3306 with a blank root password. Import both .sql files from earlier.
12. Testing can be configured in a disposable docker env.

Before working on your local copy, please use a seperate branch.
If there are tests in the package, please make sure you add tests for your work, unless it will hit a public CCP service.