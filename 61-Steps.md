# Steps

## Server runs

```go
type Server struct {
    addr    net.Addr
    sched   *Scheduler
    lowVRAM bool
}
```

It exposes a couple of routes.

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

    r.POST("/api/me", s.WhoamiHandler)

    r.POST("/api/signout", s.SignoutHandler)
    // deprecated
    r.DELETE("/api/user/keys/:encodedKey", s.SignoutHandler)

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

## API receives a generate request

Generate request accepted and validated.

```go
// GenerateRequest describes a request sent by [Client.Generate]. While you
// have to specify the Model and Prompt fields, all the other fields have
// reasonable defaults for basic uses.
type GenerateRequest struct {
    // Model is the model name; it should be a name familiar to Ollama from
    // the library at https://ollama.com/library
    Model string `json:"model"`

    // Prompt is the textual prompt to send to the model.
    Prompt string `json:"prompt"`

    // Suffix is the text that comes after the inserted text.
    Suffix string `json:"suffix"`

    // System overrides the model's default system message/prompt.
    System string `json:"system"`

    // Template overrides the model's default prompt template.
    Template string `json:"template"`

    // Context is the context parameter returned from a previous call to
    // [Client.Generate]. It can be used to keep a short conversational memory.
    Context []int `json:"context,omitempty"`

    // Stream specifies whether the response is streaming; it is true by default.
    Stream *bool `json:"stream,omitempty"`

    // Raw set to true means that no formatting will be applied to the prompt.
    Raw bool `json:"raw,omitempty"`

    // Format specifies the format to return a response in.
    Format json.RawMessage `json:"format,omitempty"`

    // KeepAlive controls how long the model will stay loaded in memory following
    // this request.
    KeepAlive *Duration `json:"keep_alive,omitempty"`

    // Images is an optional list of raw image bytes accompanying this
    // request, for multimodal models.
    Images []ImageData `json:"images,omitempty"`

    // Options lists model-specific options. For example, temperature can be
    // set through this field, if the model supports it.
    Options map[string]any `json:"options"`

    // Think controls whether thinking/reasoning models will think before
    // responding. Can be a boolean (true/false) or a string ("high", "medium", "low")
    // for supported models. Needs to be a pointer so we can distinguish between false
    // (request that thinking _not_ be used) and unset (use the old behavior
    // before this option was introduced)
    Think *ThinkValue `json:"think,omitempty"`

    // DebugRenderOnly is a debug option that, when set to true, returns the rendered
    // template instead of calling the model.
    DebugRenderOnly bool `json:"_debug_render_only,omitempty"`
}
```

* model name is validated (only syntax), fail if invalid
* searches manifests for existing model that matches the name
* if none matches, fail
* try to parse the model spec; this is on the server side, it is mostly metadata

```go
type Model struct {
    Name           string `json:"name"`
    Config         ConfigV2
    ShortName      string
    ModelPath      string
    ParentModel    string
    AdapterPaths   []string
    ProjectorPaths []string
    System         string
    License        []string
    Digest         string
    Options        map[string]any
    Messages       []api.Message

    Template *template.Template
}
```

This spec includes a model configuration:

```go
type ConfigV2 struct {
    ModelFormat   string   `json:"model_format"`
    ModelFamily   string   `json:"model_family"`
    ModelFamilies []string `json:"model_families"`
    ModelType     string   `json:"model_type"` // shown as Parameter Size
    FileType      string   `json:"file_type"`  // shown as Quantization Level
    Renderer      string   `json:"renderer,omitempty"`
    Parser        string   `json:"parser,omitempty"`

    RemoteHost  string `json:"remote_host,omitempty"`
    RemoteModel string `json:"remote_model,omitempty"`

    // used for remotes
    Capabilities []string `json:"capabilities,omitempty"`
    ContextLen   int      `json:"context_length,omitempty"`
    EmbedLen     int      `json:"embedding_length,omitempty"`
    BaseName     string   `json:"base_name,omitempty"`

    // required by spec
    Architecture string `json:"architecture"`
    OS           string `json:"os"`
    RootFS       RootFS `json:"rootfs"`
}
```

* if the request has keepalive 0, unload a model; via `s.sched.expireRunner(m)`
* validate various parameter combinations, e.g. raw and template and system
* validate and setup up parsers (model dependent), and options, like "thinking"
* express request options as model capabilities

## Scheduling a runner (from server)

```
r, m, opts, err := s.scheduleRunner(c.Request.Context(), name.String(), caps, req.Options, req.KeepAlive)
```

This will return a `llm.LlamaServer` on success.

* validate request again (e.g. model capabilities)
* pass processing on to `s.sched.GetRunner(ctx, model, opts, keepAlive)`

## Getting a runner (server)

The scheduler will return a runnerRef.

```go
func (s *Scheduler) GetRunner(c context.Context, m *Model, opts api.Options, sessionDuration *api.Duration) (chan *runnerRef, chan error) {
    ...
}
```

We are now requesting an LLM, via `LLMRequest`

```go

type LlmRequest struct {
    ctx             context.Context //nolint:containedctx
    model           *Model
    opts            api.Options
    sessionDuration *api.Duration
    successCh       chan *runnerRef
    errCh           chan error
    schedAttempts   uint
}
```

* we do again a few more validations
* we may use a "loaded" runner, that is we do not need to do anything else
* otherwise, we send the LlmRequest to a the pendingRequests channel

The scheduler works on pendingRequests and completed requests:

```
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

The `processPending` loop works on pending requests.

* load a model via `ggml, err := llm.LoadModel(pending.model.ModelPath, 1024)`

## Loading a model on request

* model must be in ggml format; we only need the path to the file, and a

```go
func LoadModel(model string, maxArraySize int) (*ggml.GGML, error) {
    if _, err := os.Stat(model); err != nil {
        return nil, err
    }

    f, err := os.Open(model)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    ggml, err := ggml.Decode(f, maxArraySize)
    return ggml, err
}
```

The `GGML` struct is just a wrapper for the moment:

```
type GGML struct {
    container
    model
    Length int64
}
```

### GGML decode

* decode uses a private `container` type which wraps a model

```go

type container interface {
    Name() string
    Decode(io.ReadSeeker) (model, error)
}

type model interface {
    KV() KV
    Tensors() Tensors
}
```

### GGUF decode

```go
func (c *containerGGUF) Decode(rs io.ReadSeeker) (model, error) {
    if err := binary.Read(rs, c.ByteOrder, &c.Version); err != nil {
        return nil, err
    }

    var err error
    switch c.Version {
    case 1:
        err = binary.Read(rs, c.ByteOrder, &c.V1)
    case 2:
        err = binary.Read(rs, c.ByteOrder, &c.V2)
    default:
        err = binary.Read(rs, c.ByteOrder, &c.V3)
    }
    if err != nil {
        return nil, err
    }

    model := newGGUF(c)
    if err := model.Decode(rs); err != nil {
        return nil, err
    }

    return model, nil
}
```

* decode key-values (metadata)
* decode tensors (no loading, just shapes)

```
        tensor := Tensor{
            Name:   name,
            Kind:   kind,
            Offset: offset,
            Shape:  shape[:],
        }
```

### Scheduler loop (cont)

This is the model swapping part. Either we load a first model, or we are
checking if we can actually load one.

```go
                    if loadedCount == 0 {
                        // No models loaded. Load the model but prefer the best fit.
                        slog.Debug("loading first model", "model", pending.model.ModelPath)
                        s.loadFn(pending, ggml, gpus, false)
                        break
                    }

                    // More than one loaded model, so we have to see if the
                    // new one fits

                    needEvict := s.loadFn(pending, ggml, gpus, true)
                    if !needEvict {
                        slog.Debug("new model fits with existing models, loading")
                        break
                    }
```

The scheduler has a load method.

```go
func (s *Scheduler) load(req *LlmRequest, f *ggml.GGML, gpus discover.GpuInfoList, requireFull bool) bool {
```

Finally, after a few validations, we pass on to llama.Load:

```
err := llama.Load(req.ctx, gpus, requireFull)
```

There can only one model be loaded at a time, but multiple requests to the same model work.

```go
    // activeLoading is the model that we are currently working on loading,
    // including by evicting one or more other models. We can only load
    // one model at a time but new requests to models that already loaded can
    // happen in parallel
    activeLoading llm.LlamaServer
```

* if there is no "llama" server, we request a new one

```go
llama, err = s.newServerFn(gpus, req.model.ModelPath, f, req.model.AdapterPaths, req.model.ProjectorPaths, req.opts, numParallel)
```

## NewLlamaServer (server)

```
func NewLlamaServer(gpus discover.GpuInfoList, modelPath string, f *ggml.GGML, adapters, projectors []string, opts api.Options, numParallel int) (LlamaServer, error) {
```

At this point, the two underlying implementations become visible:

* OllamaEngine ("new engine")
* classic LLama

The new works via model abstraction: `textProcessor, err = model.NewTextProcessor(modelPath)`

The classic llama calls `func LoadModelFromFile(modelPath string, params ModelParams) (*Model, error) { ... }`


### Starting a runner subprocess

* setup env to use available and compatible shared libraries
* get a random unassigned port

```go
        if port == 0 {
            slog.Debug("ResolveTCPAddr failed, using random port")
            port = rand.Intn(65535-49152) + 49152 // get a random port in the ephemeral range
        }
```

There is a flag to the runner to distinguish the "new engine":

```
        if textProcessor != nil {
            // New engine
            // TODO - if we have failure to load scenarios, add logic to retry with the old runner
            params = append(params, "--ollama-engine")
        }
```

The `llmServer` struct is assembled:

```
        s := llmServer{
            port:           port,
            cmd:            exec.Command(exe, params...),
            status:         NewStatusWriter(os.Stderr),
            options:        opts,
            modelPath:      modelPath,
            loadRequest:    loadRequest,
            llamaModel:     llamaModel,
            llamaModelLock: &sync.Mutex{},
            textProcessor:  textProcessor,
            numParallel:    numParallel,
            sem:            semaphore.NewWeighted(int64(numParallel)),
            totalLayers:    f.KV().BlockCount() + 1,
            loadStart:      time.Now(),
            done:           make(chan error, 1),
        }
```

The `runner` subcommand is executed, the first param is the subcommand:

```
        params := []string{"runner"}
```

The process is started in the background with:

```
if err = s.cmd.Start(); err != nil { ... }
```


The `LlamaServer` is returned.

```
type LlamaServer interface {
    ModelPath() string
    Load(ctx context.Context, gpus discover.GpuInfoList, requireFull bool) error
    Ping(ctx context.Context) error
    WaitUntilRunning(ctx context.Context) error
    Completion(ctx context.Context, req CompletionRequest, fn func(CompletionResponse)) error
    Embedding(ctx context.Context, input string) ([]float32, error)
    Tokenize(ctx context.Context, content string) ([]int, error)
    Detokenize(ctx context.Context, tokens []int) (string, error)
    Close() error
    VRAMSize() uint64 // Total VRAM across all GPUs
    TotalSize() uint64
    VRAMByGPU(gpuID string) uint64
    Pid() int
}
```

The two variants currently are:

```go
type llamaServer struct {
    llmServer

    ggml     *ggml.GGML
    gpus     discover.GpuInfoList // The set of GPUs covered by the memory estimate
    estimate MemoryEstimate
}

type ollamaServer struct {
    llmServer

    mem *ml.BackendMemory
}
```

## Loading a model in the runner

### Classic llama

Here, we cross the boundary to llama.cpp:

```
func LoadModelFromFile(modelPath string, params ModelParams) (*Model, error) {
    cparams := C.llama_model_default_params()
    cparams.n_gpu_layers = C.int(params.NumGpuLayers)
    cparams.main_gpu = C.int32_t(params.MainGpu)
    cparams.use_mmap = C.bool(params.UseMmap)
    cparams.vocab_only = C.bool(params.VocabOnly)

    ...

    m := Model{c: C.llama_model_load_from_file(C.CString(modelPath), cparams)}
    if m.c == nil {
        return nil, fmt.Errorf("unable to load model: %s", modelPath)
    }

    return &m, nil
```

The model is a wrapper around C:

```
type Model struct {
    c *C.struct_llama_model
}
```

### OllamaEngine

Uses the model interface:

```
// Model implements a specific model architecture, defining the forward pass and any model-specific configuration
type Model interface {
    Forward(ml.Context, input.Batch) (ml.Tensor, error)

    Backend() ml.Backend
    Config() config
}
```


