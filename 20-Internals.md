# Internals

## Overview

Ollama does a few things for convenience:

* manage **model life cycle**: download, store, customize, publish, remove
* **interactive chat**
* expose various **API** (native, openai compat)

## Inspired by docker

* ollama runs as a server, service
* cli makes request to the local (or some remote server)

## Ollama registry

The ollama project maintains a registry, which is modelled after docker. With a
bit of work, you can run your own registry (cf.
[#2388](https://github.com/ollama/ollama/issues/2388#issuecomment-1989307410)).

However, you can download files from other servers, too.

## Ollama model location

* model files reside under a dot folder (depending on platform)

```
$ tree -d /usr/share/ollama/.ollama/
/usr/share/ollama/.ollama/
└── models
    ├── blobs
    └── manifests
        ├── hf.co
        │   ├── arcee-ai
        │   │   └── SuperNova-Medius-GGUF
        │   ├── nomic-ai
        │   │   └── nomic-embed-text-v2-moe-gguf
        │   ├── TheBloke
        │   │   └── Mistral-7B-Instruct-v0.1-GGUF
        │   ├── unsloth
        │   │   ├── Nanonets-OCR-s-GGUF
        │   │   └── Qwen3-Coder-30B-A3B-Instruct-GGUF
        │   └── xtuner
        │       └── llava-llama-3-8b-v1_1-gguf
        └── registry.ollama.ai
            └── library
                ├── all-minilm
                ├── bge-large
                ├── gemma2
                ├── gemma3
                ├── gemma3n
                ├── gpt-oss
                ├── granite3.2
                ├── granite3.2-vision
                ├── granite-embedding
                ├── llama3
                ├── llama3.2
                ├── llama3.2-vision
                ├── minicpm-v
                ├── mistral
                ├── mistral-nemo
                ├── mistral-small
                ├── mistral-small3.1
                ├── mistral-small3.2
                ├── moondream
                ├── mxbai-embed-large
                ├── nomic-embed-text
                ├── nomic-embed-text-v2
                ├── paraphrase-multilingual
                ├── qwen2.5-coder
                ├── qwen2.5vl
                ├── qwen3
                ├── qwen3-coder
                ├── smollm
                └── smollm2

47 directories
```

We have manifests and blobs. Manifests link to blobs. You can download from
various model providers and also run your own. Default "registry.ollama.ai"
library contains official models (and variants people push/host there).

```
$ tree /usr/share/ollama/.ollama/models/manifests/ | head -40
/usr/share/ollama/.ollama/models/manifests/
├── hf.co
│   ├── arcee-ai
│   │   └── SuperNova-Medius-GGUF
│   │       └── Q8_0
│   ├── nomic-ai
│   │   └── nomic-embed-text-v2-moe-gguf
│   │       └── latest
│   ├── TheBloke
│   │   └── Mistral-7B-Instruct-v0.1-GGUF
│   │       └── latest
│   ├── unsloth
│   │   ├── Nanonets-OCR-s-GGUF
│   │   │   └── BF16
│   │   └── Qwen3-Coder-30B-A3B-Instruct-GGUF
│   │       └── UD-Q4_K_XL
│   └── xtuner
│       └── llava-llama-3-8b-v1_1-gguf
│           └── latest
└── registry.ollama.ai
    └── library
        ├── all-minilm
        │   └── latest
        ├── bge-large
        │   └── latest
        ├── gemma2
        │   └── 2b
        ├── gemma3
        │   ├── 12b
        │   ├── 1b
        │   ├── 270m
        │   ├── 27b
        │   ├── 27b-it-qat
        │   ├── 4b
        │   ├── 4b-it-qat
        │   └── latest
        ├── gemma3n
        │   └── latest
        ├── gpt-oss
        │   ├── 20b
```

## Manifest file

A manifest groups the files that make up a model, including the raw weights, templates, license.

```json
$ cat /usr/share/ollama/.ollama/models/manifests/registry.ollama.ai/library/gemma3n/latest | jq .
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
  "config": {
    "mediaType": "application/vnd.docker.container.image.v1+json",
    "digest": "sha256:8eac5d7750c5bd24a9c556890ecb08ff749fefb7c4952b8962c5e7835aef21be",
    "size": 491
  },
  "layers": [
    {
      "mediaType": "application/vnd.ollama.image.model",
      "digest": "sha256:38e8dcc30df4eb0e29eaf5c74ba6ce3f2cd66badad50768fc14362acfb8b8cb6",
      "size": 7547579904
    },
    {
      "mediaType": "application/vnd.ollama.image.template",
      "digest": "sha256:e0a42594d802e5d31cdc786deb4823edb8adff66094d49de8fffe976d753e348",
      "size": 358
    },
    {
      "mediaType": "application/vnd.ollama.image.license",
      "digest": "sha256:1adbfec9dcf025cbf301c072f3847527468dcfa399da7491ee4a1c9e9f1b33e9",
      "size": 8363
    }
  ]
}
```

A sampling of the media types:

```
$ find /usr/share/ollama/.ollama/models/manifests/registry.ollama.ai/library -type f | \
    xargs jq -rc '.layers[].mediaType' | sort | uniq -c | sort -nr

     45 application/vnd.ollama.image.model
     41 application/vnd.ollama.image.license
     36 application/vnd.ollama.image.params
     35 application/vnd.ollama.image.template
      8 application/vnd.ollama.image.system
      4 application/vnd.ollama.image.projector
```

### MediaType: model

See: [22-GGUF.md](22-GGUF.md)

### MediaType: licence

```
$ find /usr/share/ollama/.ollama/models/manifests/registry.ollama.ai/library -type f | \
    xargs jq -rc '.layers[] | select(.mediaType == "application/vnd.ollama.image.license") | .digest'  | shuf -n 10

sha256:832dd9e00a68dd83b3c3fb9f5588dad7dcf337a0db50f7d9483f310cd292e92e
sha256:c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4
sha256:43070e2d4e532684de521b885f385d0841030efa2b1a20bafb76133a5e1379c1
sha256:dd084c7d92a3c1c14cc09ae77153b903fd2024b64a100a0cc8ec9316063d2dbc
sha256:43070e2d4e532684de521b885f385d0841030efa2b1a20bafb76133a5e1379c1
sha256:cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30
sha256:3e2c24001f9ef57bf7ec959a3658fbb49cdad113cdf394c264da9d16f9bdd132
sha256:43070e2d4e532684de521b885f385d0841030efa2b1a20bafb76133a5e1379c1
sha256:a406579cd136771c705c521db86ca7d60a6f3de7c9b5460e6193a2df27861bde
sha256:2e68075caee4c571d43cb9d8d636415d8fe0fe59f695370cf56bc9f872f1ff3f
```

Example license:

```
$ cat /usr/share/ollama/.ollama/models/blobs/sha256-832dd9e00a68dd83b3c3fb9f5588dad7dcf337a0db50f7d9483f310cd292e92e | head

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
```


### MediaType: params

```
$ find /usr/share/ollama/.ollama/models/manifests/registry.ollama.ai/library -type f | \
    xargs jq -rc '.layers[] | select(.mediaType == "application/vnd.ollama.image.params") | .digest'  | shuf -n 10

sha256:832dd9e00a68dd83b3c3fb9f5588dad7dcf337a0db50f7d9483f310cd292e92e
sha256:c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4
sha256:43070e2d4e532684de521b885f385d0841030efa2b1a20bafb76133a5e1379c1
sha256:dd084c7d92a3c1c14cc09ae77153b903fd2024b64a100a0cc8ec9316063d2dbc
sha256:43070e2d4e532684de521b885f385d0841030efa2b1a20bafb76133a5e1379c1
sha256:cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30
sha256:3e2c24001f9ef57bf7ec959a3658fbb49cdad113cdf394c264da9d16f9bdd132
sha256:43070e2d4e532684de521b885f385d0841030efa2b1a20bafb76133a5e1379c1
sha256:a406579cd136771c705c521db86ca7d60a6f3de7c9b5460e6193a2df27861bde
sha256:2e68075caee4c571d43cb9d8d636415d8fe0fe59f695370cf56bc9f872f1ff3f
```

Example params:

```
$ cat /usr/share/ollama/.ollama/models/blobs/sha256-339e884a40f6708bc761d367f0c08e448d5bb6f16b3961c340e44e0e4835a004 | jq .
{
  "stop": [
    "<end_of_turn>"
  ],
  "top_k": 64,
  "top_p": 0.95
}

$ cat /usr/share/ollama/.ollama/models/blobs/sha256-e0daf17ff83eace4813f9e8554b262f6cc33ad880ff8df41a156ff9ef5522ddb | jq .
{
  "temperature": 0.15
}
```

Typical parameters:

```
$ ollama ls  | grep -v NAME | awk '{print $1}' | \
    xargs -I {} ollama show --parameters {} | sort  | awk '{print $1}' | uniq -c | sort -nr
     44 stop
     24 temperature
     15 top_p
     12 top_k
     12
      7 num_ctx
      5 repeat_penalty
      2 min_p
      1 num_keep
```



### MediaType: system

A system prompt.

```shell
$ cat /usr/share/ollama/.ollama/models/blobs/sha256-70a4dab5e1d14953cc95c7d4ee1003f05c2474a39cd07fef2f8975c776455d33

You are Mistral Small 3.1, a Large Language Model (LLM) created by Mistral AI,
a French startup headquartered in Paris.  You power an AI assistant called Le
Chat.  Your knowledge base was last updated on 2023-10-01.

When you're not sure about some information, you say that you don't have the
information and don't make up anything.  If the user's question is not clear,
ambiguous, or does not provide enough context for you to accurately answer the
question, you do not try to answer it right away and you rather ask the user to
clarify their request (e.g. "What are some good restaurants around me?" =>
"Where are you?" or "When is the next flight to Tokyo" => "Where do you travel
from?").  You are always very attentive to dates, in particular you try to
resolve dates (e.g. "yesterday" is {yesterday}) and when asked about
information at specific dates, you discard information that is at another date.
You follow these instructions in all languages, and always respond to the user
in the language they use or request.  Next sections describe the capabilities
that you have.

# WEB BROWSING INSTRUCTIONS

You cannot perform any web search or access internet to open URLs, links etc.
If it seems like the user is expecting you to do so, you clarify the situation
and ask the user to copy paste the text directly in the chat.

# MULTI-MODAL INSTRUCTIONS

You have the ability to read images, but you cannot generate images. You also
cannot transcribe audio files or videos.
```

Or just:

```
$ cat /usr/share/ollama/.ollama/models/blobs/sha256-75357d685f238b6afd7738be9786fdafde641eb6ca9a3be7471939715a68a4de
You are a helpful assistant.
```

### MediaType: projector



### ollama show

Display model properties and parameters with show subcommand:

```
$ ollama show gemma3:4b-it-qat
  Model
    architecture        gemma3
    parameters          4.3B
    context length      131072
    embedding length    2560
    quantization        Q4_0

  Capabilities
    completion
    vision

  Parameters
    top_p          0.95
    stop           "<end_of_turn>"
    temperature    1
    top_k          64

```


## Ollama Codebase

About 78K LOC Go. Ollama vendors [ggml](https://github.com/ggml-org/ggml) and
[llama.cpp](https://github.com/ggml-org/llama.cpp), which are C and C++ projects.

```
$ tokei
===============================================================================
 Language            Files        Lines         Code     Comments       Blanks
===============================================================================
 GNU Style Assembly      1            6            6            0            0
 Autoconf                1            4            4            0            0
 C                      12        26550        21425          947         4178
 C Header               87       122062       103618        11551         6893
 CMake                   8         1739         1403          112          224
 C++                    58        95684        75391         5559        14734
 C++ Header              3        26105        18285         4684         3136
 CSS                     1           34           29            0            5
 Dockerfile              1          123          104            3           16
 Go                    327        78674        64719         4230         9725
 HTML                    1            9            9            0            0
 JavaScript              2           13           12            1            0
 JSON                   38       147592       147592            0            0
 Objective-C             2         6816         5624          233          959
 PowerShell              2          265          242            9           14
 Protocol Buffers        1          333           98          178           57
 Shell                   7          608          476           39           93
 SVG                     1            9            9            0            0
 Plain Text              4       190212            0       167280        22932
 TSX                     2          129          122            0            7
 TypeScript              9          494          407           14           73
-------------------------------------------------------------------------------
 Markdown               25         3659            0         2439         1220
 |- BASH                 1            1            1            0            0
 |- Dockerfile           3           22           21            0            1
 |- Go                   1           23           22            0            1
 |- INI                  2           18           16            0            2
 |- JavaScript           1           43           34            1            8
 |- JSON                 1          504          504            0            0
 |- PowerShell           3            5            5            0            0
 |- Python               2           96           78            2           16
 |- Shell               15          629          624            0            5
 |- TypeScript           1           18           15            0            3
 (Total)                           5018         1320         2442         1256
===============================================================================
 Total                 593       701120       439575       197279        64266
===============================================================================
```


## High Level Request Flow

1. **Request Reception**: Client sends HTTP request to Ollama server (port
   11434) via CLI, REST API, or OpenAI-compatible endpoint

2. **Request Routing**: Server routes request through `server/routes.go` to
   appropriate handler (`/api/generate`, `/api/chat`, etc.)

3. **Model Management**: Scheduler (`server/sched.go`) checks if model is
   loaded; if not, loads model into memory with appropriate backend (CUDA/ROCm/Metal)

4. **Prompt Processing**: Request is parsed and prepared for inference by the
   runner system

5. **Inference Execution**: Specialized runner (`runner/`) performs model
   inference using optimized backends for token generation

6. **Response Streaming**: Generated tokens are streamed back to client via
   HTTP streaming (default mode) or returned as complete response

7. **Resource Management**: System manages GPU/CPU resources and model
   lifecycle for subsequent requests


## Model files
