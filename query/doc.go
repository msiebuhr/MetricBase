/*
Package Query (partially) implements support for graphite queries.

The sub-package queryParser generates a parse-tree from a query-string and this
packge converts that query to an executable tree.

This package has major TODO-stuff:

 * Naming of nodes vs. ast. vs. Source vs. Functions is crap. Despite multiple
   attempts, I can't come up with better names.
 * The package organization is ... non-existant it's all in a big pile of code.
 * Not extensible in any way (i.e. not possible to add new functions)

*/
package query
