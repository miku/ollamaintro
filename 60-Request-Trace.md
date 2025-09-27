# Tracing a request

Outline of the request flow for a completion request to the ollama api (coming
from some client).

## Request handling

* completion request arrives

## Runner instance

## Tokenization

## Batching



## Code Snippets

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
