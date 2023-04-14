install-dep:
	go install github.com/goreleaser/goreleaser@latest
	go install github.com/caarlos0/svu@latest

release:
	git tag $(shell svu next)
	git push --tags
	goreleaser release --clean