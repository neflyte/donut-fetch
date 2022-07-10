donut-fetch:
	CGO_ENABLED=0 go build -ldflags '-s -w' ./cmd/donut-fetch
.PHONY: donut-fetch

clean-cache:
	if [ -d "$(HOME)/.cache/donut-fetch" ]; then rm -Rf "$(HOME)/.cache/donut-fetch"; fi
.PHONY: clean-cache

clean-state:
	if [ -f "$(HOME)/.config/donut-fetch/state.json" ]; then rm -f "$(HOME)/.config/donut-fetch/state.json"; fi
.PHONY: clean-state

clean:
	rm ./donut-fetch || true
.PHONY: clean
