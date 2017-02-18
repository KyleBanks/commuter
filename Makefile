# Builds the commuter CLI.
build:
	@go build github.com/KyleBanks/commuter/cmd/commuter 
.PHONY: build

# Runs an example commuter request for travel duration.
example: | build
	@./commuter -to "Toronto, Canada"
.PHONY: example

# Runs test suit, vet, golint, and fmt.
sanity:
	@echo "---------------- TEST ----------------"
	@go list ./... | grep -v vendor/ | xargs go test --cover 

	@echo "---------------- VET ----------------"
	@go list ./... | grep -v vendor/ | xargs go vet 

	@echo "---------------- LINT ----------------"
	@go list ./... | grep -v vendor/ | xargs golint

	@echo "---------------- FMT ----------------"
	@go list ./... | grep -v vendor/ | xargs go fmt
.PHONY: sanity
