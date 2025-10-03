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

As of 2025, ollama support a number of model capabilities:

```go
package model

type Capability string

const (
    CapabilityCompletion = Capability("completion")
    CapabilityTools      = Capability("tools")
    CapabilityInsert     = Capability("insert")
    CapabilityVision     = Capability("vision")
    CapabilityEmbedding  = Capability("embedding")
    CapabilityThinking   = Capability("thinking")
)

func (c Capability) String() string {
    return string(c)
}
```

* completion
* tools
* insert
* vision
* embedding
* thinking



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

* support for ollama on the site: [https://huggingface.co/docs/hub/en/ollama](https://huggingface.co/docs/hub/en/ollama)
* format: `hf.co/{username}/{repository}`

```
$ ollama run hf.co/DevQuasar/swiss-ai.Apertus-8B-Instruct-2509-GGUF:Q4_K_M
```

> By default, the `Q4_K_M` quantization scheme is used, when it’s present inside the model repo.

### Custom Server

* you can use the API to register blobs and then name a model that references the blobs

## Quantization

* model size depends on number of parameters, typically millions, billions, up to trillions
* one technique to reduce model size (other: lora, pruning, knowledge
  destillation, parameter sharing; cf. [model shrinking
techniques](https://web.dev/articles/llm-sizes#model-shrinking))
* quantization results in an quantization error, the less the better

> Quantization: Reducing the precision of weights from floating-point numbers
> (such as, 32-bit) to lower-bit representations (such as, 8-bit).

Common data types:

* FP32 (full precision)
* FP16 (half precision)
* BF16 (a shortened 16-bit version of 32-bit floats)

> It preserves the approximate dynamic range of 32-bit floating-point numbers
> by retaining 8 exponent bits, but supports only an 8-bit precision rather
> than the 24-bit significand of the binary32 format. -- [bfloat16 floating-point format](https://en.wikipedia.org/wiki/Bfloat16_floating-point_format)

> Our results show that deep learning training using BFLOAT16 tensors achieves
> the same state-of-the-art (SOTA) results across domains as FP32 tensors in
> the same number of iterations and with no changes to hyper-parameters. -- [A Study of BFLOAT16 for Deep Learning Training](https://arxiv.org/pdf/1905.12322)

* INT8 (1 byte)

### Suffixes

* https://huggingface.co/docs/hub/en/gguf#quantization-types
* todo: https://medium.com/@paul.ilvez/demystifying-llm-quantization-suffixes-what-q4-k-m-q8-0-and-q6-k-really-mean-0ec2770f17d3

Deep dive: [A Visual Guide to Quantization](https://www.maartengrootendorst.com/blog/quantization/)

## Using the API

* todo: curl
* todo: parameters
* todo: using the API from Go, SDK

## What are the different model types?

* [Model Types](42-Model-Types.md)


## Working with a text model

### Single completion

### Chat completion


## Working with an image model

### Passing an image


## Working with an embedding model

## Request Trace

* [Request Trace](60-Request-Trace.md)

## Tokenization

* character tokens
* word tokens
* subword tokens

### BPE in Go

* todo: BPE in go; cf. https://eli.thegreenplace.net/2024/tokens-for-llms-byte-pair-encoding-in-go/
* [x/bpe](x/bpe)

## What happens in the server?

Server exposes chat and management API.

```go
    // General
    r.HEAD("/", func(c *gin.Context) { c.String(http.StatusOK, "Ollama is running") })
    r.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "Ollama is running") })
    r.HEAD("/api/version", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"version": version.Version}) })
    r.GET("/api/version", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"version": version.Version}) })

    // Local model cache management (new implementation is at end of function)
    r.POST("/api/pull", s.PullHandler)
    r.POST("/api/push", s.PushHandler)
    r.HEAD("/api/tags", s.ListHandler)
    r.GET("/api/tags", s.ListHandler)
    r.POST("/api/show", s.ShowHandler)
    r.DELETE("/api/delete", s.DeleteHandler)

    // Create
    r.POST("/api/create", s.CreateHandler)
    r.POST("/api/blobs/:digest", s.CreateBlobHandler)
    r.HEAD("/api/blobs/:digest", s.HeadBlobHandler)
    r.POST("/api/copy", s.CopyHandler)

    // Inference
    r.GET("/api/ps", s.PsHandler)
    r.POST("/api/generate", s.GenerateHandler)
    r.POST("/api/chat", s.ChatHandler)
    r.POST("/api/embed", s.EmbedHandler)
    r.POST("/api/embeddings", s.EmbeddingsHandler)

    // Inference (OpenAI compatibility)
    r.POST("/v1/chat/completions", openai.ChatMiddleware(), s.ChatHandler)
    r.POST("/v1/completions", openai.CompletionsMiddleware(), s.GenerateHandler)
    r.POST("/v1/embeddings", openai.EmbeddingsMiddleware(), s.EmbedHandler)
    r.GET("/v1/models", openai.ListMiddleware(), s.ListHandler)
    r.GET("/v1/models/:model", openai.RetrieveMiddleware(), s.ShowHandler)
```

* server manages runners
* server can load and load runners
* server can have keep multiple models running at the same time
* it will swap out models, as requests require

A few more tasks:

* checking model capabilities ("images.go"), matched against model properties, read from gguf

```go
// CheckCapabilities checks if the model has the specified capabilities returning an error describing
// any missing or unknown capabilities
func (m *Model) CheckCapabilities(want ...model.Capability) error {
    available := m.Capabilities()
    var errs []error

    // Map capabilities to their corresponding error
    capToErr := map[model.Capability]error{
        model.CapabilityCompletion: errCapabilityCompletion,
        model.CapabilityTools:      errCapabilityTools,
        model.CapabilityInsert:     errCapabilityInsert,
        model.CapabilityVision:     errCapabilityVision,
        model.CapabilityEmbedding:  errCapabilityEmbedding,
        model.CapabilityThinking:   errCapabilityThinking,
    }

    for _, cap := range want {
        err, ok := capToErr[cap]
        if !ok {
            slog.Error("unknown capability", "capability", cap)
            return fmt.Errorf("unknown capability: %s", cap)
        }

        if !slices.Contains(available, cap) {
            errs = append(errs, err)
        }
    }
    ...
```

### Scheduler

The scheduler schedules models, load and unloads them, keeps track of available memory.

```go
// Returns immediately, spawns go routines for the scheduler which will shutdown when ctx is done
func (s *Scheduler) Run(ctx context.Context) {
        slog.Debug("starting llm scheduler")
        go func() {
                s.processPending(ctx)
        }()

        go func() {
                s.processCompleted(ctx)
        }()
}
```


## What happens in the runner?

The runner is a slim server accepting the actual completion request.

Ollama has two runner options.

```
    if newRunner {
        return ollamarunner.Execute(args)
    } else {
        return llamarunner.Execute(args)
    }
```

The llamarunner is the classic runner, ollamarunner is a new, Go abstraction
with a default backend using GGML.

### A pure Go abstraction

* `ml/backend.go` contains a pure Go abstraction of the execution runtime
* a `Tensor` interface
* a `Context` interface, containing

Currently one backend based on ggml. The `New` constructor is less than 400
SLOC, after which the `ml.Backend` abstraction can be used.

```go
func NewBackend(modelPath string, params BackendParams) (Backend, error) {
    if backend, ok := backends["ggml"]; ok {
        return backend(modelPath, params)
    }

    return nil, fmt.Errorf("unsupported backend")
}
```

Setting up the ggml based backend:

* read gguf file, parse metadata
* collect available backends (cpu, gpu, accelarators)
* calculate required memory across devices
* assign each layer to a backend

```
    layers := make([]deviceBufferType, blocks)
    for i := range layers {
        layers[i] = assignLayer(i)
    }
```

* allocate tensors (with name mappings)

```go
createTensor := func(t tensor, bts []C.ggml_backend_buffer_type_t, layer int) *C.struct_ggml_tensor
```

* assemble backend struct and return

### Layers to GPU

Layers are assigned to multiple GPUs, if available.

```
// GPULayers is a set of layers to be allocated on a single GPU
type GPULayers struct {
    // ID is the identifier of the GPU, as reported in DeviceMemory
    ID string

    // Layers is a set of layer indicies to load
    Layers []int
}
```



## Where and how are the model files stored?

### Ollama model location

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

### Manifest file

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


## What is a GGUF file?

* [GGUF](46-GGUF.md)

## Customizations

* prompt
* request options, e.g. structured output


## Model Customization Options

* Modelfile, https://ollama.readthedocs.io/en/modelfile/
* change parameters: top-k, top-n, system prompt, adapter


## Model Evaluation

How do I know, how a model performs a certain task?

> Evaluations

* there is a large number of benchmarks and benchmark datasets
* the open llm leaderboard hosted on HF has been discontinued in 2025
* benchmarks will test a benchmark dataset (in isolation), limited insight into general usefulness

> [Misrepresented Technological Solutions in Imagined Futures: The Origins and
> Dangers of AI Hype in the Research
> Community](https://arxiv.org/pdf/2408.15244) (2024)

* LLM are multitask learners and hence an LLM is not really made for a single purpose


