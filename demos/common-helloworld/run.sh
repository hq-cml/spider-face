#!/bin/sh
cd /data/share/golang/src/github.com/hq-cml/spider-face/demos/helloworld/
go build ./
if [ $? -ne 0 ]; then
    echo "failed"
    exit
fi
#mv helloworld /data/share/golang/src/github.com/hq-cml/spider-face/demos/helloworld/
/data/share/golang/src/github.com/hq-cml/spider-face/demos/helloworld/helloworld
