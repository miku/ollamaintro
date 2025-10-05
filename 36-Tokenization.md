# Tokenization

* character tokens
* word tokens
* subword tokens

## BPE outline

Introduced by [Neural Machine Translation of Rare Words with Subword Units](https://aclanthology.org/P16-1162/) (2016) for machine translation.

> Previous work addresses the translation of out-of-vocabulary words by backing
> off to a dictionary. In this paper, we introduce a simpler and more ef-
> fective approach, making the NMT model capable of open-vocabulary translation
> by encoding rare and unknown words as se- quences of subword units

Algorithm outline.

* tokenize text into words
* start with initial vocabulary of all characters that make up the words ("a", "b", س  ...)
* find most common pair
* substitute pair of chars with a single token, e.g. ["e", "r"] -> ["er"]
* repeat until you arrive at a vocabulary size of choice

E.g. gemma3 has 262144 tokens; llama3.2 128253, ...

```
$ ollama show --verbose gemma3 | grep "tokenizer.ggml.tokens"
    tokenizer.ggml.tokens                         [<pad> ...+262144 more]

```

## BPE in Go

* BPE in go; cf. https://eli.thegreenplace.net/2024/tokens-for-llms-byte-pair-encoding-in-go/

