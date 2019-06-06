#!/bin/sh
cd /data/share/golang/src/github.com/hq-cml/spider-face/demos/common-html/
go build -o helloworld ./
if [ $? -ne 0 ]; then
    echo "failed"
    exit
fi
/data/share/golang/src/github.com/hq-cml/spider-face/demos/common-html/helloworld
