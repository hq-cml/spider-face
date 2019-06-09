#!/bin/sh
cd /data/share/golang/src/github.com/hq-cml/spider-face/demos/spider-ui/
go build -o ui ./
if [ $? -ne 0 ]; then
    echo "failed"
    exit
fi
/data/share/golang/src/github.com/hq-cml/spider-face/demos/spider-ui/ui
