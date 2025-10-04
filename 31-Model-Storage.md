# Where and how are the model files stored?

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
