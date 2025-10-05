SHELL := /bin/bash

word2vec.model: make_model.py corpus.file
	time python make_model.py

corpus.file:
	cat *txt > corpus.file

.PHONY: clean
clean:
	rm -f corpus.file
