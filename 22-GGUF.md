# GGUF

## Overview

See also: [GGUF](https://github.com/ggml-org/ggml/blob/master/docs/gguf.md)

> GGUF is a file format for storing models for inference with GGML and
> executors based on GGML. GGUF is a binary format that is designed for fast
> loading and saving of models, and for ease of reading. Models are
> traditionally developed using PyTorch or another framework, and then
> converted to GGUF for use in GGML.

Some format evolution:

> GGML, GGMF, GGJT, and now GGUF — and it is now at **version three** of the format.

## Example: Reading a gguf file

Random, example project:
[SmolDocling-256M-preview-GGUF](https://huggingface.co/Mungert/SmolDocling-256M-preview-GGUF),
[files and
versions](https://huggingface.co/Mungert/SmolDocling-256M-preview-GGUF/tree/main)

```
$ curl -sLo SmolDocling-256M-preview-q8_0.gguf "https://huggingface.co/Mungert/SmolDocling-256M-preview-GGUF/resolve/main/SmolDocling-256M-preview-q8_0.gguf?download=true"
```

