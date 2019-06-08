#!/bin/sh
cd /data/share/golang/src/github.com/hq-cml/spider-face/demos/forum/
go build -o forum ./
if [ $? -ne 0 ]; then
    echo "failed"
    exit
fi
/data/share/golang/src/github.com/hq-cml/spider-face/demos/forum/forum
