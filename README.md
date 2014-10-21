webcp
=====

Copy a subtree of a web site, with smart page filters. Distinctive features:

- Crawl pages without saving them, in order to discover links to the pages you really want.
- Automatically crawl from the Internet Archive instead of the live site (please be nice to the Archive!)

__WORK IN PROGRESS:__ This program is not yet ready to use.

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

    webcp <url> .

By default, the crawl will fetch all linked pages up to a depth of 5, and will delay 5 seconds between subsequent requests to the same domain.

If you have a large crawl that you might need to kill and later resume, you can do that by providing a resume file:

    webcp --resume=links.txt <url> .

Planned Enhancements
--------------------

- __Crawl__ linked pages to the maximum depth, but only __save__ pages whose URLs/MIME types match certain filters.
- Stay on the same domain, or set of domains.
- Crawl from the Internet Wayback Machine instead of from the live site, with fancy date filtering to get the page version you want.
