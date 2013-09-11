MetricBase
==========

For play/experimentation/single-server version of
[Graphite](https://github.com/graphite-project).

Structure
---------

Front-ends are pieces of code that talk to the world around it. They, in turn,
talk to a backend.

Back-ends are chained, and commands modified/passed from one to the next until the
request is processed in it's entirety.

Building
--------

Regardless of whether it's used or not, `libleveldb` has to be installed in
some form or other for the server to compile. On Debian/Ubuntu it's
`libleveldb-dev`, on OS X w. Homebrew it's `leveldb`.

    go get github.com/msiebuhr/MetricBase
	go install github.com/msiebuhr/MetricBase/metricBaseServer

Start the server

    ./metricBaseServer

It listens for Graphites text protocol on TCP port 2003 and has a webserver
running on http://localhost:8080/ that serves a simple front-end from
`/http-pub/`.
