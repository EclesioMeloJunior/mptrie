.PHONY: test
test:
	@go test -timeout 30s ./... 

cover:
	@go test -coverprofile mptcoverage.html ./... 
	@go tool cover -html=./mptcoverage.html && unlink mptcoverage.html