OLDGOPATH=$GOPATH
export GOPATH=`pwd`

go build

export GOPATH=$OLDGOPATH
