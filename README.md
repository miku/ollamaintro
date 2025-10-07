# What happens when you type a prompt into Ollama?

> Intro to Ollama Workshop at [GoLab](https://golab.io) 2025, 2025-10-05,
> [Martin Czygan](https://de.linkedin.com/in/martin-czygan-58348842)

## Overview

Ollama is a popular tool to run LLM (and multimodal and embedding models) on
your own machine. As of 09/2025 it has over 150K stars on GitHub.

In 2025, there are numerous other tools to run models locally:

* [llama.cpp](https://github.com/ggml-org/llama.cpp) (wrapped by ollama)
* [llamafile](https://github.com/Mozilla-Ocho/llamafile)
* [vLLM](https://github.com/vllm-project/vllm) (used by various cloud services)
* [OpenLLM](https://github.com/bentoml/OpenLLM)
* [Lemonade](https://lemonade-server.ai/)
* [Jan.ai](https://github.com/menloresearch/jan)
* [yzma](https://github.com/hybridgroup/yzma)
* and more ...

And even more user interfaces of various kinds.

As of 09/2025, of the [25809 repositories](https://github.com/topics/llm) on GitHub
tagged [llm], ollama seems to be among the top ten.

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

