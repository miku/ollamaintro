# Ollama Intro Workshop

> Workshop at [GoLab](https://golab.io) 2025, 2025-10-05, [Martin
> Czygan](https://de.linkedin.com/in/martin-czygan-58348842)

## Intro and setup

* [Intro](10-Intro.md)
* [Motivation](15-Motivation.md)

## LLM background and milestones

* [Background](20-Background.md)

## ollama cli

* cli talks with the server

```sh
$ ollama
Usage:
  ollama [flags]
  ollama [command]

Available Commands:
  serve       Start ollama
  create      Create a model
  show        Show information for a model
  run         Run a model
  stop        Stop a running model
  pull        Pull a model from a registry
  push        Push a model to a registry
  signin      Sign in to ollama.com
  signout     Sign out from ollama.com
  list        List models
  ps          List running models
  cp          Copy a model
  rm          Remove a model
  help        Help about any command

Flags:
  -h, --help      help for ollama
  -v, --version   Show version information

Use "ollama [command] --help" for more information about a command.
```

## Running a model

```
$ ollama run llama3:latest
```

* download model from default registry to local cache
* drops you into a chat interface

[![](static/screenshot-2025-09-30-174533-intel-n150-alder-lake-llama3.png)](https://github.com/miku/ollamaintro/blob/main/static/ollama-chat-n150-llama3.gif?raw=true)

You can pass a prompt directly as an argument.


## What models?

* ollama library lists models in the default registry:
  [library](https://ollama.com/library); also [list in
README](https://github.com/ollama/ollama/?tab=readme-ov-file#model-library)
* between 815MB ("gemma3:1b") and 404GB ("deepseek-r1:671b") in size

> You should have at least 8 GB of RAM available to run the 7B models, 16 GB to
> run the 13B models, and 32 GB to run the 33B models. --
> [README.md](https://github.com/ollama/ollama/)

* HF support ollama: [https://huggingface.co/docs/hub/en/ollama](https://huggingface.co/docs/hub/en/ollama)


### Registry

* model name will get resolved to the registry manifest

[![](static/screenshot-2025-10-01-000836-ollama-direct-downloader.png)](https://ollama-direct-downloader.vercel.app/)

```sh
$ curl -sL "https://registry.ollama.ai/v2/library/gemma3/manifests/latest" | jq .
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
  "config": {
    "mediaType": "application/vnd.docker.container.image.v1+json",
    "digest": "sha256:b6ae5839783f2ba248e65e4b960ab15f9c4b7118db285827dba6cba9754759e2",
    "size": 489
  },
  "layers": [
    {
      "mediaType": "application/vnd.ollama.image.model",
      "digest": "sha256:aeda25e63ebd698fab8638ffb778e68bed908b960d39d0becc650fa981609d25",
      "size": 3338792448
    },
    {
      "mediaType": "application/vnd.ollama.image.template",
      "digest": "sha256:e0a42594d802e5d31cdc786deb4823edb8adff66094d49de8fffe976d753e348",
      "size": 358
    },
    {
      "mediaType": "application/vnd.ollama.image.license",
      "digest": "sha256:dd084c7d92a3c1c14cc09ae77153b903fd2024b64a100a0cc8ec9316063d2dbc",
      "size": 8432
    },
    {
      "mediaType": "application/vnd.ollama.image.params",
      "digest": "sha256:3116c52250752e00dd06b16382e952bd33c34fd79fc4fe3a5d2c77cf7de1b14b",
      "size": 77
    }
  ]
}
```

### Huggingface

* todo: download link
* todo: custom model manifest on own server

## Quantization

* model size depends on number of parameters, typically millions, billions, up to trillions
* one technique to reduce model size (other: lora, pruning, knowledge
  destillation, parameter sharing; cf. [model shrinking
techniques](https://web.dev/articles/llm-sizes#model-shrinking))

> Quantization: Reducing the precision of weights from floating-point numbers
> (such as, 32-bit) to lower-bit representations (such as, 8-bit).

## Using the API

* todo: curl
* todo: parameters

## What are the different model types?

* todo: custom modes; maybe gifcities, example for clip


