# Intro

## Quick poll

* [ ] Who has already used Ollama?
* [ ] Who has run an LLM locally?
* [ ] Who has worked with model APIs before?

## Workshop goals

* [ ] Better understanding how Ollama works under the hood
* [ ] Learn about Model lifecycle and customization
* [ ] Build practical applications with embeddings in Go
* [ ] Create a visual search tool for handwritten notes

## Local and remote machine setup

* Install ollama on your laptop, you can follow the instructions here: [https://ollama.com/download](https://ollama.com/download)
* We will have access to a cloud GPU server instance, running a RTX 4000 SFF GPU with 20GB VRAM (an efficient GPU, that consumes at most 70W)

Most examples will work well on Linux and MacOS, but Windows WSL should also
work.

## Checkpoint: installation

```
$ ollama --version
$ ollama ls
$ curl localhost:11434
```

## Workshop material

All material available on
[github.com/miku/ollamaintro](https://github.com/miku/ollamaintro). Please
clone, as if contains also code examples and project scaffolding.


