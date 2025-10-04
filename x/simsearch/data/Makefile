SHELL := /bin/bash

pg1661.json: pg1661.py pg1661.txt
	cat pg1661.txt | python pg1661.py > pg1661.json

pg1661_sentences.json: pg1661.txt
	sentences < pg1661.txt | jq -rc -R '[., inputs] | to_entries | map({id: .key, text: .value}) | .[]' > pg1661_sentences.json
