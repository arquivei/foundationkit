.PHONY: send send-%

send:
	$(eval URL := http://localhost:8686/ping/v1)
	@if [ -f "$(R)" ] ; then \
		cat $(R) | curl -v -H 'Content-Type: application/json' -d @- "$(URL)"; \
	else \
		echo "Invalid request file '$(R)'" > /dev/stderr; \
	fi

send-without-trace: R ?= docs/requests/3-5000.json
send-without-trace: send

run:
	GOOGLE_APPLICATION_CREDENTIALS=gcp-credentials.json go run cmd/api/*.go 