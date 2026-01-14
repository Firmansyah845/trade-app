run:
	@go run . server

run-linter:
	rm -f report.html
	golangci-lint run ./...

run-mega-linter:
	rm -rf megalinter-reports
	npx mega-linter-runner --flavor go