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
