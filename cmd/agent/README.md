**build agent with custom flags**

go build -ldflags "-X main.buildVersion=v0.0.20 -X main.buildDate=04-10-2024 -X main.buildCommit=some_commit" -o agent cmd/agent/*.go