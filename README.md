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

	go get github.com/msiebuhr/MetricBase
	go build ./bin/MetricBase/

Start the server

	./MetricBase

It listens for the Graphite text protocol on TCP port 2003 and has a webserver
running on http://localhost:8080/ that serves a simple front-end from
`/http-pub/`.
