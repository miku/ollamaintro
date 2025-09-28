# Tracing a request

Outline of the request flow for a completion request to the ollama api (coming
from some client).

## Request handling

* completion request arrives

## Runner instance

## Tokenization

## Batching



## Code Snippets

### Mostly llamarunner

```go
type CompletionRequest struct {
    Prompt  string
    Format  json.RawMessage
    Images  []ImageData
    Options *api.Options

    Grammar string // set before sending the request to the subprocess
}
```

A sequence is an abstraction over one request.

```go
type Sequence struct {
    // batch index
    iBatch int

    // number of tokens predicted so far
    numPredicted int

    // prompt inputs left to evaluate
    inputs []input

    // inputs that have been added to a batch but not yet submitted to Decode
    pendingInputs []input

    // tokens that have been generated but not returned yet (e.g. for stop sequences)
    pendingResponses []string

    // input cache being used by this sequence
    cache *InputCacheSlot

    // channel to send responses over
    responses chan string

    // channel to stop decoding (such as if the remote connection is closed)
    quit chan bool

    // number of tokens to predict
    numPredict int

    samplingCtx *llama.SamplingContext

    // channel to send back the embedding if embedding only
    embedding chan []float32

    // stop sequences
    stop []string

    // number of inputs to keep at the beginning when shifting context window
    numKeep int

    // true if an embedding are to be returned instead of text generation
    embeddingOnly bool

    doneReason llm.DoneReason

    // Metrics
    startProcessingTime time.Time
    startGenerationTime time.Time
    numDecoded          int
    numPromptInputs     int
}
```

The input can be a token or an image:

```
// input is an element of the prompt to process, either
// a token or an image embedding (generated from a vision projector)
type input struct {
    token int

    // embed is an image embedding
    embed []float32
}
```

The  server, on loading a model, allocates a maximum number of requests to
handle in parallel.

```go
        // runner/llamarunner/runner.go:807
        s.seqs = make([]*Sequence, s.parallel)
```

In llamarunner, the server keeps a number of slots for the sequences. The server will wait, until the a slow becomes available. Once it is available

```go
    // Ensure there is a place to put the sequence, released when removed from s.seqs
    if err := s.seqsSem.Acquire(r.Context(), 1); err != nil {
        if errors.Is(err, context.Canceled) {
            slog.Info("aborting completion request due to client closing the connection")
        } else {
            http.Error(w, fmt.Sprintf("Failed to acquire semaphore: %v", err), http.StatusInternalServerError)
        }
        return
    }

    s.mu.Lock()
    found := false
    for i, sq := range s.seqs {
        if sq == nil {
            seq.cache, seq.inputs, err = s.cache.LoadCacheSlot(seq.inputs, true)
            if err != nil {
                s.mu.Unlock()
                s.seqsSem.Release(1)
                http.Error(w, fmt.Sprintf("Failed to load cache: %v", err), http.StatusInternalServerError)
                return
            }

            s.seqs[i] = seq
            s.cond.Signal()
            found = true
            break
        }
    }
    s.mu.Unlock()
```

The server runs in a loop:

```
    for {
        select {
        case <-ctx.Done():
            return
        default:
            err := s.processBatch(tokenBatch, embedBatch)
            if err != nil {
                panic(err)
            }

            tokenBatch.Clear()
            embedBatch.Clear()
        }
    }
```

The server has a processBatch function; called from server run, which is started in a background thread.

```
    for range s.seqs {
        seqIdx = (seqIdx + 1) % len(s.seqs)
        seq := s.seqs[seqIdx]

        if seq == nil {
            continue
        }
```

In this processBatch function, we call `llama_decode` through the wrapper (cf.
[llama_context::decode](https://github.com/ggml-org/llama.cpp/blob/e6d65fb02d553bd79cad94e517cdca18b687788d/src/llama-context.cpp#L958-L1260),
which itself calls
[[llama_context::process_ubatch](https://github.com/ggml-org/llama.cpp/blob/e6d65fb02d553bd79cad94e517cdca18b687788d/src/llama-context.cpp#L732-L794))

There, we build the graph, set the inputs and run `graph_compute`, cf. [llama-context.cpp#L755-L784](https://github.com/ggml-org/llama.cpp/blob/e6d65fb02d553bd79cad94e517cdca18b687788d/src/llama-context.cpp#L775-L784)

```cpp
    // set the input data for the input tensors
    {
        //const auto t_start_us = ggml_time_us();

        res->set_inputs(&ubatch);

        //LLAMA_LOG_INFO("graph set inputs time: %.3f ms\n", (ggml_time_us() - t_start_us)/1000.0);
    }

    const auto status = graph_compute(res->get_gf(), ubatch.n_tokens > 1);
```


### New backend architecture

* new abstraction layer for backends (cf. `ml/backend.go`):

```go
type Backend interface {
    // Close frees all memory associated with this backend
    Close()

    Load(ctx context.Context, progress func(float32)) error

    // BackendMemory returns the memory allocations that were made for this model
    BackendMemory() BackendMemory

    Config() fs.Config
    Get(name string) Tensor
    NewContext() Context
    NewContextSize(size int) Context
}
```
