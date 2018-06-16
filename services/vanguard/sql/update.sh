#!/bin/bash
wget -q -O - https://www.fuzzwork.co.uk/dump/mysql-latest.tar.bz2 | \
    tar -xjO | sed 's/InnoDB/TokuDB/' | gzip -9c > eve.gz