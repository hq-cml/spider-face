#!/bin/sh
cd /data/share/golang/src/github.com/hq-cml/spider-face/demos/quick-html/
go build -o quick ./
if [ $? -ne 0 ]; then
    echo "failed"
    exit
fi
/data/share/golang/src/github.com/hq-cml/spider-face/demos/quick-html/quick
