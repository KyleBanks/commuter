# Builds and installs the commuter CLI.
install:
	@go install github.com/KyleBanks/commuter 
.PHONY: install

# Runs an example commuter request for travel duration.
example: | install
	@commuter -to "Toronto, Canada"
.PHONY: example

include github.com/KyleBanks/make/go/sanity
include github.com/KyleBanks/make/go/release
include github.com/KyleBanks/make/misc/precommit
