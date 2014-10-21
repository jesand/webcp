webcp
=====

Copy a subtree of a web site, with smart URL filters

Compiling
---------

    go get -t github.com/jesand/webcp
    go test github.com/jesand/webcp/...
    go install github.com/jesand/webcp

Usage Examples
--------------

See the usage notes:

    webcp -h

Download a URL and its sub-pages into the current directory:

    webcp <url>

By default, the crawl will fetch all linked pages up to a depth of 5, and will delay 5 seconds between subsequent requests to the same domain.

If you have a large crawl that you might need to kill and later resume, you can do that by providing a resume file:

    webcp <url> --resume=links.txt

Planned Enhancements
--------------------

- __Crawl__ linked pages to the maximum depth, but only __save__ pages whose URLs/MIME types match certain filters.
- Stay on the same domain, or set of domains.
- Crawl from the Internet Wayback Machine instead of from the live site, with fancy date filtering to get the page version you want.
