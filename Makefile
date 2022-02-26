BAND = Weedrams
SETLIST = Blocks

.PHONY: help clean generate-spick
.DEFAULT_GOAL := help

clean: ## clean all output files
	rm -rf out

generate-spick: ## generate spick pdf
	./script/generate-spick.sh "$(BAND)" "$(SETLIST)"

generate-setlist: ## generate Setlist
	@mkdir -p out
	@cp ~/data/obsidian/notes/music/howlers/Howlers\ Repertoire.md Howlers/Repertoire.md
	@go run main.go list --band "$(BAND)" --gig  "$(SETLIST)"

kindle-email: ## send to kindle via email
	./script/send.sh "$(BAND)" "$(SETLIST)"

kindle: ## copy to plugged in kindle
	cp "out/$(BAND)@$(SETLIST).pdf" /Volumes/Kindle/documents

all-setlist:
	find $(BAND)/songs -type f -name "*.odt" -printf "%f\n" | sed -e 's/\.odt$$//' | sed -e 's/^/* /' > $(BAND)/gigs/all.md

all-songs: SETLIST=all
all-songs: all-setlist generate-spick

all: ## generate all (howlers & weedrams)
	$(MAKE) SETLIST=all BAND=Howlers generate-spick forscore
	$(MAKE) SETLIST=all BAND=Weedrams generate-spick forscore

weedrams-blocks: ## generate weedrams blocks
	$(MAKE) SETLIST=Blocks BAND=Weedrams generate-spick

forscore: ## update cloud storage for forScore
	cp -p $(BAND)/songs/*.pdf ~/Library/Mobile\ Documents/com\~apple\~CloudDocs/forScore/
	cp -p out/*.pdf ~/Library/Mobile\ Documents/com\~apple\~CloudDocs/forScore/

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
