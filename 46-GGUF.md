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

![](static/313174776-c3623641-3a1d-408e-bfaf-1b7c4e16aa63.png)

## Highlights

* Single-file deployment: they can be easily distributed and loaded, and do not require any external files for additional information.
* Extensible: new features can be added to GGML-based executors/new information can be added to GGUF models without breaking compatibility with existing models.
* mmap compatibility: models can be loaded using mmap for fast loading and saving.
* Easy to use: models can be easily loaded and saved using a small amount of code, with no need for external libraries, regardless of the language used.
* Full information: all information needed to load a model is contained in the model file, and no additional information needs to be provided by the user.

