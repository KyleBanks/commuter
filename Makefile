# Build and run the commuter CLI.
commuter:
	@go build github.com/KyleBanks/commuter/cmd/commuter
	@./commuter -to "Toronto, Canada"
.PHONY: commuter

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
