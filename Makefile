# Build and run the commuter CLI.
commuter:
	@go build github.com/KyleBanks/commuter/cmd/commuter 
.PHONY: commuter

# Runs an example commuter request for travel duration.
example: | commuter
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
