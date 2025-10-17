# What happens when you type a prompt into Ollama?

> Intro to Ollama Workshop at [GoLab](https://golab.io) 2025, 2025-10-05,
> [Martin Czygan](https://de.linkedin.com/in/martin-czygan-58348842)

We have [Slides](Slides.md) and [examples](x/), which are mostly self contained.

## Overview

Ollama is a popular tool to run LLM (and multimodal and embedding models) on
your own machine. As of 09/2025 it has over [150K stars on GitHub](https://github.com/ollama/ollama/).

In 2025, there are numerous other tools to run models locally:

* [llama.cpp](https://github.com/ggml-org/llama.cpp) (wrapped by ollama)
* [llamafile](https://github.com/Mozilla-Ocho/llamafile) (single file distribution)
* [vLLM](https://github.com/vllm-project/vllm) (used by various cloud services, #1 for prod)
* [Docker Model Runner](https://www.docker.com/blog/introducing-docker-model-runner/) (cf. [How OCI Artifacts Will Drive Future AI Use Cases](https://www.cncf.io/blog/2025/08/27/how-oci-artifacts-will-drive-future-ai-use-cases/)))
* [OpenLLM](https://github.com/bentoml/OpenLLM)
* [Lemonade](https://lemonade-server.ai/) (AMD)
* [Jan.ai](https://github.com/menloresearch/jan)
* [yzma](https://github.com/hybridgroup/yzma) (lightweight, ffi, Go)
* [ggml](https://github.com/ggml-org/ggml) (used by ollama)
* and more ...

Even more user interfaces of various kinds exist.

As of 09/2025, of the [25809 repositories](https://github.com/topics/llm) on GitHub
tagged [llm], ollama seems to be among the top ten.

Note: ollama is both open source
([MIT](https://github.com/ollama/ollama/?tab=MIT-1-ov-file#readme) licensed) and [VC
funded](https://www.ycombinator.com/companies/ollama).

## Before the workshop / quick start checklist

* [ ] please [install ollama](https://ollama.com/download) on your laptop (ok, if it only has a cpu)

After installation, please run the following commands to download a few models
files onto your laptop (order of preference):

```
ollama pull embeddinggemma
ollama pull llama3.2
ollama pull gemma3:270m
ollama pull gemma3
ollama pull qwen2.5vl
```

Warning: These file models may occupy **over 10GB** of disk space.


## Let's go

* [Slides](Slides.md)
