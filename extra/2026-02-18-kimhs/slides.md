# What Happens When You Type a Prompt into Ollama?

> A deep dive into the internals of local LLM inference

---

## Agenda

1. Ollama Architecture Overview
2. The Request Journey (prompt → tokens → response)
3. Multimodal: How Vision Models Work
4. GGUF: The Model Format
5. Build Your Own: Creating a Toy Architecture

---

## Ollama Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        API Layer                             │
│   HTTP :11434  │  REST + OpenAI compat  │  gin-gonic        │
├─────────────────────────────────────────────────────────────┤
│                       Scheduler                              │
│   Model lifecycle  │  Request queuing  │  Resource mgmt     │
├─────────────────────────────────────────────────────────────┤
│                       Runner Layer                           │
│   llamarunner (llama.cpp)  │  ollamarunner (ggml direct)    │
└─────────────────────────────────────────────────────────────┘
```

- One subprocess per loaded model
- Docker-inspired UX: `pull`, `run`, `push`
- Content-addressable storage (like git)

---

## The Request Journey

```
User types: "Why is the sky blue?"
     │
     ▼
┌──────────────────────────────────────────────────────────────┐
│ 1. API receives GenerateRequest                              │
│ 2. Model lookup and validation                               │
│ 3. Template rendering (prompt assembly)                      │
│ 4. Tokenization                                              │
│ 5. Scheduler allocates sequence slot                         │
│ 6. Runner processes batch                                    │
│ 7. Token sampling loop                                       │
│ 8. Streaming response back                                   │
└──────────────────────────────────────────────────────────────┘
```

---

## Step 1: API Request

```go
// github.com/ollama/ollama/api/types.go
type GenerateRequest struct {
    Model    string         `json:"model"`
    Prompt   string         `json:"prompt"`
    Images   []ImageData    `json:"images,omitempty"`  // multimodal
    Stream   *bool          `json:"stream,omitempty"`
    Options  map[string]any `json:"options"`
    Think    *ThinkValue    `json:"think,omitempty"`   // reasoning
}
```

Simple curl example:

```bash
curl localhost:11434/api/generate -d '{
  "model": "gemma3:270m",
  "prompt": "Why is the sky blue?"
}'
```

---

## Step 2: Model Lookup

```go
// Server looks up model from manifests
type Model struct {
    Name           string
    ModelPath      string         // path to GGUF blob
    ProjectorPaths []string       // vision projector (multimodal)
    Template       *template.Template
    Options        map[string]any
}
```

Model storage follows Docker pattern:

```
~/.ollama/models/
├── blobs/                    # GGUF files (sha256-...)
└── manifests/
    └── registry.ollama.ai/
        └── library/
            └── gemma3/
                └── 270m      # JSON manifest
```

---

## Step 3: Template Rendering

Templates transform raw prompt into model-specific format:

```
<start_of_turn>user
Why is the sky blue?<end_of_turn>
<start_of_turn>model
```

Debug with `_debug_render_only`:

```bash
curl localhost:11434/api/generate -d '{
  "model": "gemma3:270m",
  "prompt": "Why is the sky blue?",
  "_debug_render_only": true
}'
```

---

## Step 4: Tokenization

```go
// Text → token IDs using BPE
"The sky is blue" → [464, 6766, 318, 4171]
```

Tokenizers are model-specific (stored in GGUF metadata):
- BPE (Byte-Pair Encoding) - most common
- SentencePiece - Google models
- WordPiece - BERT-style

---

## Step 5: Scheduler & Sequence Slots

```go
type Sequence struct {
    inputs          []input        // tokens or image embeddings
    cache           *InputCacheSlot // KV cache for context
    responses       chan string    // streaming output
    samplingCtx     *llama.SamplingContext
    stop            []string       // stop sequences
}

// input: either a token or image embedding
type input struct {
    token int
    embed []float32
}
```

Scheduler manages:
- Model loading/unloading (5min keep-alive default)
- Parallel request slots
- GPU/CPU memory allocation

---

## Step 6: Runner Forward Pass

```go
// ollamarunner main loop
func (s *Server) run(ctx context.Context) {
    for {
        activeBatch, err = s.forwardBatch(activeBatch)
        s.computeBatch(activeBatch)  // calls ggml
    }
}
```

The batch is processed through:
1. Token embeddings lookup
2. Transformer blocks (attention + FFN)
3. Final normalization
4. Output projection → logits

---

## Step 7: Token Sampling

```go
// Sample next token from logits
vocabSize := len(outputs) / batchSize
token, err := seq.sampler.Sample(
    outputs[i*vocabSize : (i+1)*vocabSize])

// Check for end of sequence
if s.model.(model.TextProcessor).Is(token, model.SpecialEOS) {
    s.removeSequence(i, llm.DoneReasonStop)
    continue
}
```

Sampling parameters: temperature, top_k, top_p, repeat_penalty

---

## Step 8: Streaming Response

```go
// Token → text, streamed back to client
func respFunc(resp api.GenerateResponse) error {
    fmt.Print(resp.Response)  // each token as it's generated
    return nil
}

client.Generate(ctx, req, respFunc)
```

Response includes metadata:
```json
{
  "response": "The sky appears blue because...",
  "done": true,
  "total_duration": 1234567890,
  "prompt_eval_count": 8,
  "eval_count": 42
}
```

---

## Multimodal: How Vision Works

Vision models combine image encoder + text decoder:

```
┌──────────────────────────────────────────────────────────┐
│  Image bytes → Vision Encoder → Embeddings              │
│       ↓                                                  │
│  [img-0] token placeholder in prompt                    │
│       ↓                                                  │
│  Combined: image_embeddings + text_tokens               │
│       ↓                                                  │
│  Transformer processes both modalities                  │
│       ↓                                                  │
│  Text output: "A cat sitting on a keyboard..."          │
└──────────────────────────────────────────────────────────┘
```

---

## Multimodal: Code Example

```go
// x/sdk-image/main.go
client, _ := api.ClientFromEnvironment()
imgData, _ := os.ReadFile("photo.jpg")

req := &api.GenerateRequest{
    Model:  "qwen2.5vl",
    Prompt: "Describe this image in detail",
    Images: []api.ImageData{imgData},
}

client.Generate(ctx, req, func(resp api.GenerateResponse) error {
    fmt.Print(resp.Response)
    return nil
})
```

Models with vision: llava, qwen2.5vl, moondream, llama3.2-vision

---

## Multimodal: Image Processing

```go
// inputs processes prompt and images into inputs
func (s *Server) inputs(prompt string, images []llm.ImageData) (
    []*input.Input, []ml.Context, multimodalStore, error) {

    // Split prompt on [img-N] tags
    // Tokenize text segments
    // Decode images through vision encoder
    // Interleave embeddings
}
```

Template with image placeholder:
```
<|im_start|>user
[img-0]
What is in this picture?<|im_end|>
<|im_start|>assistant
```

---

## GGUF: The Model Format

**GGUF** = GGML Universal Format

```
┌─────────────────────────────────────────────────┐
│  Header: magic, version, tensor count, kv count │
├─────────────────────────────────────────────────┤
│  Key-Value Pairs (metadata)                     │
│  - general.architecture = "llama"               │
│  - llama.context_length = 4096                  │
│  - tokenizer.ggml.model = "gpt2"                │
├─────────────────────────────────────────────────┤
│  Tensor Infos (name, shape, type, offset)       │
├─────────────────────────────────────────────────┤
│  Tensor Data (weights, mmap-ready)              │
└─────────────────────────────────────────────────┘
```

- Single-file deployment
- mmap for fast loading
- Quantization support (Q4_K_M, Q8_0, etc.)

---

## GGUF: Reading Metadata

```go
// cmd/readgguf/main.go
f, err := gguf.Open(modelPath)
defer f.Close()

fmt.Printf("Magic: %s\n", f.Magic)
fmt.Printf("Version: %d\n", f.Version)
fmt.Printf("Tensors: %d\n", f.NumTensors())

for _, kv := range f.KeyValues() {
    fmt.Printf("%s = %v\n", kv.Key, kv.Value)
}

for _, t := range f.TensorInfos() {
    fmt.Printf("%s: %v [%s]\n", t.Name, t.Shape, t.Type)
}
```

---

## GGUF: Tensor Structure

Common tensor naming pattern:

```
token_embd.weight          [vocab_size, hidden_dim]
blk.0.attn_q.weight        [hidden_dim, hidden_dim]
blk.0.attn_k.weight        [hidden_dim, kv_dim]
blk.0.attn_v.weight        [hidden_dim, kv_dim]
blk.0.attn_output.weight   [hidden_dim, hidden_dim]
blk.0.ffn_up.weight        [hidden_dim, ffn_dim]
blk.0.ffn_down.weight      [ffn_dim, hidden_dim]
blk.0.attn_norm.weight     [hidden_dim]
blk.0.ffn_norm.weight      [hidden_dim]
...
output_norm.weight         [hidden_dim]
output.weight              [hidden_dim, vocab_size]
```

---

## Build Your Own: Toy Architecture

Goal: Create a minimal model with random weights that Ollama can load

Steps:
1. Define architecture (PyTorch or direct GGUF)
2. Generate/save weights as GGUF
3. Set required metadata
4. Create manifest and register with Ollama

---

## Toy Model: PyTorch Definition

```python
# example_toymodel.py
import torch
import torch.nn as nn

class TinyLM(nn.Module):
    def __init__(self, vocab_size=256, hidden_dim=64, n_layers=2):
        super().__init__()
        self.token_embd = nn.Embedding(vocab_size, hidden_dim)

        self.blocks = nn.ModuleList([
            nn.ModuleDict({
                'attn': nn.Linear(hidden_dim, hidden_dim),
                'ffn': nn.Linear(hidden_dim, hidden_dim),
                'norm': nn.LayerNorm(hidden_dim),
            }) for _ in range(n_layers)
        ])

        self.output_norm = nn.LayerNorm(hidden_dim)
        self.lm_head = nn.Linear(hidden_dim, vocab_size)

    def forward(self, x):
        h = self.token_embd(x)
        for blk in self.blocks:
            h = h + blk['attn'](h)
            h = blk['norm'](h)
            h = h + blk['ffn'](h)
        h = self.output_norm(h)
        return self.lm_head(h)
```

---

## Toy Model: GGUF Writer

```python
# example_gguf_writer.py
import gguf
import numpy as np

writer = gguf.GGUFWriter("toymodel.gguf", "llama")

# Required metadata
writer.add_architecture()
writer.add_context_length(512)
writer.add_embedding_length(64)
writer.add_block_count(2)
writer.add_head_count(4)
writer.add_vocab_size(256)

# Add tensors with random weights
vocab, hidden = 256, 64
writer.add_tensor("token_embd.weight",
    np.random.randn(vocab, hidden).astype(np.float32))

for i in range(2):
    writer.add_tensor(f"blk.{i}.attn_q.weight",
        np.random.randn(hidden, hidden).astype(np.float32))
    # ... add other block tensors

writer.add_tensor("output.weight",
    np.random.randn(vocab, hidden).astype(np.float32))

writer.write_header_to_file()
writer.write_kv_data_to_file()
writer.write_tensors_to_file()
writer.close()
```

---

## Toy Model: Creating the Manifest

```bash
# Create a Modelfile that references the GGUF
cat > Modelfile << 'EOF'
FROM ./toymodel.gguf
PARAMETER temperature 0.7
PARAMETER num_ctx 512
EOF

# Register with Ollama
ollama create toymodel -f Modelfile

# Verify it's available
ollama list | grep toymodel
```

---

## Toy Model: Testing

```bash
# Test with ollama runner directly
ollama runner --ollama-engine --model ./toymodel.gguf --port 8080 &

# Load the model (specify cache parameters)
curl -X POST localhost:8080/load -d '{
  "Operation": 2, "Parallel": 1,
  "BatchSize": 512, "KvSize": 2048
}'

# Generate tokens
curl -X POST localhost:8080/completion \
  -d '{"prompt": "Hello", "n_predict": 10}'
# Output: random tokens like <0x1F37><0x9B3>... (untrained!)
```

This proves the GGUF structure is valid and the architecture works.

---

## Complete Request Flow Summary

```
curl → API (validate) → Scheduler (load model)
                              ↓
                        Runner subprocess
                              ↓
              ┌───────────────────────────────┐
              │ 1. Template render            │
              │ 2. Tokenize prompt            │
              │ 3. Forward pass (ggml/llama)  │
              │ 4. Sample token               │
              │ 5. Check stop conditions      │
              │ 6. Stream response            │
              │ 7. Loop until done            │
              └───────────────────────────────┘
                              ↓
                    Response streamed back
```

---

## Key Takeaways

1. **Three layers**: API → Scheduler → Runner (subprocess)
2. **Templates** transform prompts for each model family
3. **GGUF** is the universal model format (single file, mmap-ready)
4. **Multimodal** = vision encoder embeddings + text tokens
5. **Custom architectures** need proper GGUF metadata + tensor names
6. **Token sampling** loop with temperature, top_k, top_p

---

## Resources

- Ollama source: https://github.com/ollama/ollama
- GGUF spec: https://github.com/ggml-org/ggml/blob/master/docs/gguf.md
- llama.cpp: https://github.com/ggml-org/llama.cpp
- This workshop: https://github.com/miku/ollamaintro

---

## Questions?

```
      ██████  ██      ██       █████  ███    ███  █████
     ██    ██ ██      ██      ██   ██ ████  ████ ██   ██
     ██    ██ ██      ██      ███████ ██ ████ ██ ███████
     ██    ██ ██      ██      ██   ██ ██  ██  ██ ██   ██
      ██████  ███████ ███████ ██   ██ ██      ██ ██   ██
```

