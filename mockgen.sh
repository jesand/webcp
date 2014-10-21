#!/bin/bash
#
# This script generates mock classes using the mockgen program
# installed with gomock. See:
#
# http://godoc.org/code.google.com/p/gomock/gomock
#
# To install gomock, run:
#
# go get code.google.com/p/gomock/gomock
# go install code.google.com/p/gomock/mockgen
#

# Generate mock classes
mockgen -package="crawl" -destination="crawl/mock_queue_test.go" github.com/jesand/webcp/crawl CrawlQueueStorage
