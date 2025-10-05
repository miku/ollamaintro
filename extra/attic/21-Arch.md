# Ollama architecture

## Core ideas

* local execution (of models in GGUF format)
* model lifecycle (registry, customizations)
* various interfaces: chat, API

It is possible to get gguf files from any remote

## Client/Server

* server running as a service, exposing various API
* cli as API client
* easy to run or share a remote instance


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


## Path of a generate request

> Generate a response for a given prompt with a provided model. -- [docs](https://ollama.readthedocs.io/en/api/#generate-a-completion)


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

The completion model parses a model struct, hovering over the manifest.

```go
func GetModel(name string) (*Model, error) {
    mp := ParseModelPath(name)
    manifest, digest, err := GetManifest(mp)
    if err != nil {
        return nil, err
    }

    model := &Model{
        Name:      mp.GetFullTagname(),
        ShortName: mp.GetShortTagname(),
        Digest:    digest,
        Template:  template.DefaultTemplate,
    }

    if manifest.Config.Digest != "" {
        filename, err := GetBlobsPath(manifest.Config.Digest)
        if err != nil {
            return nil, err
        }

        configFile, err := os.Open(filename)
        if err != nil {
            return nil, err
        }
        defer configFile.Close()

        if err := json.NewDecoder(configFile).Decode(&model.Config); err != nil {
            return nil, err
        }
    }

    for _, layer := range manifest.Layers {
        filename, err := GetBlobsPath(layer.Digest)
        if err != nil {
            return nil, err
        }

        switch layer.MediaType {
        case "application/vnd.ollama.image.model":
            model.ModelPath = filename
            model.ParentModel = layer.From
        case "application/vnd.ollama.image.embed":
            // Deprecated in versions  > 0.1.2
            // TODO: remove this warning in a future version
            slog.Info("WARNING: model contains embeddings, but embeddings in modelfiles have been deprecated and will be ignored.")
        case "application/vnd.ollama.image.adapter":
            model.AdapterPaths = append(model.AdapterPaths, filename)
        case "application/vnd.ollama.image.projector":
            model.ProjectorPaths = append(model.ProjectorPaths, filename)
        case "application/vnd.ollama.image.prompt",
            "application/vnd.ollama.image.template":
            bts, err := os.ReadFile(filename)
            if err != nil {
                return nil, err
            }

            model.Template, err = template.Parse(string(bts))
            if err != nil {
                return nil, err
            }
        case "application/vnd.ollama.image.system":
            bts, err := os.ReadFile(filename)
            if err != nil {
                return nil, err
            }

            model.System = string(bts)
        case "application/vnd.ollama.image.params":
            params, err := os.Open(filename)
            if err != nil {
                return nil, err
            }
            defer params.Close()

            // parse model options parameters into a map so that we can see which fields have been specified explicitly
            if err = json.NewDecoder(params).Decode(&model.Options); err != nil {
                return nil, err
            }
        case "application/vnd.ollama.image.messages":
            msgs, err := os.Open(filename)
            if err != nil {
                return nil, err
            }
            defer msgs.Close()

            if err = json.NewDecoder(msgs).Decode(&model.Messages); err != nil {
                return nil, err
            }
        case "application/vnd.ollama.image.license":
            bts, err := os.ReadFile(filename)
            if err != nil {
                return nil, err
            }
            model.License = append(model.License, string(bts))
        }
    }

    return model, nil
}
```

This is more like model metadata:

```

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

A first look into the request prompt, if it is empty we unload the model.

Some models have custom reponse renderers, like gptoss
[harmony](https://cookbook.openai.com/articles/openai-harmony).

Once the request input is validated, as runner is requested:

```
r, m, opts, err := s.scheduleRunner(c.Request.Context(), name.String(), caps, req.Options, req.KeepAlive)
```

Schedule runner will validated model capabilities, and will call `GetRunner`

```
    runnerCh, errCh := s.sched.GetRunner(ctx, model, opts, keepAlive)
```

The scheduler determines if the request could be satisfied by a loaded runner,
otherwise it is send to a pending requests queue.

When the model is available, we proceed to add images and preparing the prompt.

```go
  images := make([]llm.ImageData, len(req.Images))
    for i := range req.Images {
        images[i] = llm.ImageData{ID: i, Data: req.Images[i]}
    }
```

The server ask LlamaServer for completion.

```
        if err := r.Completion(c.Request.Context(), llm.CompletionRequest{
            Prompt:     prompt,
            Images:     images,
            Format:     req.Format,
            Options:    opts,
            UseHarmony: useHarmony,
        }, func(cr llm.CompletionResponse) {
            res := api.GenerateResponse{
                Model:     req.Model,
                CreatedAt: time.Now().UTC(),
                Response:  cr.Content,
                Done:      cr.Done,
                Thinking:  cr.Thinking,
                ToolCalls: cr.ToolCalls,
                Metrics: api.Metrics{
                    PromptEvalCount:    cr.PromptEvalCount,
                    PromptEvalDuration: cr.PromptEvalDuration,
                    EvalCount:          cr.EvalCount,
                    EvalDuration:       cr.EvalDuration,
                },
            }
```


## The runner process

Ollama has two runner options.

```
    if newRunner {
        return ollamarunner.Execute(args)
    } else {
        return llamarunner.Execute(args)
    }
```

## Input

```
// Input represents one token in the input stream
type Input struct {
    // Token is a single element of text.
    Token int32

    // Multimodal is represents a non-text element such as an
    // image (or part of one if the image can be processed in pieces).
    // It may be used either together with Token or on its own.
    Multimodal []Multimodal

    // MultimodalHash is a unique representation of the data
    // stored in Multimodal, used for caching and comparing
    // equality.
    MultimodalHash uint64

    // SameBatch forces the following number of tokens to be processed
    // in a single batch, breaking and extending batches as needed.
    // Useful for things like images that must be processed in one
    // shot.
    SameBatch int
}
```

## Batch

```
// Batch contains the inputs for a model forward pass
type Batch struct {
    // Inputs is the input tokens, including placeholders for multimodal inputs.
    Inputs ml.Tensor

    // Outputs are the set of indicies into Inputs for which output data should
    // be returned.
    Outputs ml.Tensor

    // Positions is the position for each Input, relative to its sequence. Equal
    // in length to Inputs.
    Positions []int32

    // Sequences is the sequence for each Input. Equal in length to Inputs.
    Sequences []int

    // Multimodal is a set of multimodal embeddings previously created by
    // EncodeMultimodal, along with an index into Inputs. Unused for text-only
    // models or for batches without multimodal elements.
    Multimodal []MultimodalIndex
}
```


## Loading and serving models

* models need to be loaded into memory
* diverse backends: CPU, GPU
* manage multiple models and requests (swapping models if required)


