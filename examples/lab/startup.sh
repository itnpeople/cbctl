#!/bin/bash
docker run -d --rm -p 1024:1024 --name cb-spider -v "$(pwd)/data/spider:/data" -e CBSTORE_ROOT=/data cloudbaristaorg/cb-spider:0.5.0
docker run -d --rm -p 1323:1323 --name cb-tumblebug --link cb-spider:cb-spider -v "$(pwd)/data/tumblebug/conf:/app/conf" -v "$(pwd)/data/tumblebug/meta_db:/app/meta_db/dat" -v "$(pwd)/data/tumblebug/log:/app/log" cloudbaristaorg/cb-tumblebug:0.5.0
docker run -d --rm -p 1470:1470 --name cb-mcks --link cb-spider:cb-spider --link cb-tumblebug:cb-tumblebug -v "$(pwd)/data/mcks:/data" -e SPIDER_URL=http://cb-spider:1024/spider -e TUMBLEBUG_URL=http://cb-tumblebug:1323/tumblebug -e CBSTORE_ROOT=/data cloudbaristaorg/cb-mcks:latest
docker ps