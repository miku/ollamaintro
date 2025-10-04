# More on models

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


Not supported: audio, video

## Modalities

* completion is the basic text interaction
* tools; [x/toolcall](x/toolcall)
* insert (adds a suffix, e.g. code completion)
* vision: joint text and image models, captioning, OCR, ...
* embedding: maps data into a dense vector space
* thinking ...

Test-time compute scaling; more tokens, better responses; related to prompt
engineering, CoT - mostly empiric. Thinking token, delimiter for "<think>" ...;
can be turned on off. Cf. deepseek-r1 ("[Incentivizing Reasoning Capability in
LLMs via Reinforcement Learning](https://arxiv.org/abs/2501.12948)" - cited by
[4110](https://scholar.google.com/scholar?cluster=2469397274690356930)); "trained on question/answer pairs"

