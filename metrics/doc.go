/*
Package metrics implement basic operations on Graphite-style metrics.

A metric has a dot-separated name, a floating-point value and a UNIX-timestamp,
ex:

    serverX.systemY.subsystemZ.users.10m_avg 523.1 1400609412
*/
package metrics
