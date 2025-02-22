
orchestrator:
	go build -o orchestrator cmd/orchestrator/main.go

agent:
	go build -o agent cmd/agent/main.go

clean:
	rm -f orchestrator
	rm -f agent

.PHONY: clean orchestrator agent
