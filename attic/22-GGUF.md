# GGUF

## Overview

See also: [GGUF](https://github.com/ggml-org/ggml/blob/master/docs/gguf.md)

> GGUF is a file format for storing **models for inference with GGML** and
> executors based on GGML. GGUF is a binary format that is designed for fast
> loading and saving of models, and for ease of reading. Models are
> traditionally developed using PyTorch or another framework, and then
> converted to GGUF for use in GGML.

Some format evolution:

> GGML, GGMF, GGJT, and now GGUF — and it is now at **version three** of the format.

Full description: [gguf.md](https://github.com/ggml-org/ggml/blob/master/docs/gguf.md)

## Format overview

* magic
* version
* number of tensors
* number of key value pairs (metadata)

![](static/313174776-c3623641-3a1d-408e-bfaf-1b7c4e16aa63.png)

## Naming convention

```
<BaseName><SizeLabel><FineTune><Version><Encoding><Type><Shard>.gguf
```

Examples:

```
Mixtral-8x7B-v0.1-KQ2.gguf
Hermes-2-Pro-Llama-3-8B-F16.gguf
Grok-100B-v1.0-Q4_0-00003-of-00009.gguf
```

## readgguf

* [x/ggufopen/](x/ggufopen)

Example cli application: readgguf. You can point it to a file or URL to get
some information about gguf file.

Random, example project:
[SmolDocling-256M-preview-GGUF](https://huggingface.co/Mungert/SmolDocling-256M-preview-GGUF),
[files and
versions](https://huggingface.co/Mungert/SmolDocling-256M-preview-GGUF/tree/main)

```
$ curl -sLo SmolDocling-256M-preview-q8_0.gguf "https://huggingface.co/Mungert/SmolDocling-256M-preview-GGUF/resolve/main/SmolDocling-256M-preview-q8_0.gguf?download=true"
```

## Quantization

* [Blind testing different quants](https://github.com/ggml-org/llama.cpp/discussions/5962)
