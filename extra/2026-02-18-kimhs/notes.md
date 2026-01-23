# Speaker Notes

## Slide: Title

- Welcome, today we'll trace exactly what happens when you type a prompt into Ollama
- This is a condensed version of a 3-hour workshop, focused on the code paths
- We'll see real code from the Ollama project, not just diagrams

## Slide: Agenda

- Start with architecture overview to set context
- Then trace a single request from API to response
- Special section on multimodal (vision models)
- GGUF format deep dive - it's the foundation
- End with hands-on: creating your own toy model

## Slide: Ollama Architecture

- **API Layer**: Standard HTTP server using gin-gonic framework
  - Port 11434 is the default
  - OpenAI-compatible endpoints make it a drop-in replacement
  - Routes: `/api/generate`, `/api/chat`, `/api/embed`, etc.

- **Scheduler**: Think of it as a model manager
  - Handles loading/unloading models based on memory
  - Default keep-alive is 5 minutes
  - Manages request queuing when model is loading

- **Runner Layer**: This is where inference happens
  - One subprocess per loaded model (can see with `pstree` or `ollama ps`)
  - Two implementations: classic llama.cpp wrapper vs newer direct ggml
  - The `--ollama-engine` flag switches between them

Key point: Docker-inspired design makes model management familiar to developers

## Slide: The Request Journey

Walk through the 8 steps briefly:
1. API receives the HTTP request
2. Validates model exists in local cache
3. Applies model-specific template to prompt
4. Converts text to token IDs
5. Scheduler ensures model is loaded, allocates slot
6. Runner processes tokens through transformer layers
7. Samples next token from probability distribution
8. Streams each token back as it's generated

Each step we'll look at actual code.

## Slide: Step 1: API Request

- `GenerateRequest` is the main struct for text generation
- Note the optional fields:
  - `Images` for multimodal - we'll cover this later
  - `Think` for reasoning models like DeepSeek-R1
  - `Stream` defaults to true - token-by-token response

Demo opportunity: Show the curl command live

## Slide: Step 2: Model Lookup

- Models are stored in a Docker-like registry structure
- `blobs/` contains actual GGUF files, named by SHA256 hash
- `manifests/` contains JSON files linking to blobs
- This enables content-addressable storage:
  - Deduplication (same weights shared between model versions)
  - Integrity verification
  - Easy caching

The `Model` struct holds everything needed to serve a model.

## Slide: Step 3: Template Rendering

- Each model family has its own template
- Templates use Go's text/template syntax
- They add special tokens that the model was trained on:
  - `<start_of_turn>`, `<end_of_turn>` for Gemma
  - `<|im_start|>`, `<|im_end|>` for Qwen
  - `[INST]`, `[/INST]` for Llama

Demo: Use `_debug_render_only` to see the actual prompt sent to model

## Slide: Step 4: Tokenization

- Text must become numbers for the neural network
- BPE (Byte-Pair Encoding) is most common
- Tokenizer trained alongside model - stored in GGUF metadata
- Typical vocab sizes: 32k-128k tokens

Fun fact: "The" and " The" are usually different tokens

## Slide: Step 5: Scheduler & Sequence Slots

- `Sequence` represents one ongoing generation
- Key fields:
  - `inputs`: queue of tokens/embeddings to process
  - `cache`: KV cache slot for this conversation (enables context)
  - `responses`: channel for streaming output

- The `input` union type is interesting:
  - Either a token ID (int) or an embedding vector (float32 slice)
  - This is how multimodal works - images become embeddings

## Slide: Step 6: Runner Forward Pass

- The main loop is surprisingly simple
- Continuously builds batches and computes them
- Batch = multiple sequences processed together for efficiency
- Each `forwardBatch` runs the full transformer:
  1. Token embeddings lookup
  2. For each layer: attention → add&norm → FFN → add&norm
  3. Final normalization
  4. Linear projection to vocab size

## Slide: Step 7: Token Sampling

- Model outputs logits (unnormalized probabilities) for each token in vocabulary
- Sampling strategy determines which token to pick:
  - Temperature: controls randomness (higher = more random)
  - Top-k: only consider k most likely tokens
  - Top-p (nucleus): consider tokens until cumulative probability reaches p

- EOS (End of Sequence) token signals completion
- Stop sequences can also trigger early termination

## Slide: Step 8: Streaming Response

- Each token is decoded and sent immediately
- Callback pattern in Go SDK
- Response includes useful metadata:
  - `total_duration`: wall-clock time
  - `prompt_eval_count`: how many tokens in prompt
  - `eval_count`: how many tokens generated

This is why you see text appearing word-by-word in chat interfaces.

## Slide: Multimodal: How Vision Works

- Vision models have two parts:
  - Image encoder (often based on CLIP/SigLIP)
  - Language model (decoder)

- The trick: images become embedding vectors that slot into the text sequence
- `[img-0]` is a placeholder that gets replaced with image embeddings
- The transformer then attends to both image and text tokens

## Slide: Multimodal: Code Example

- Simple API: just add raw image bytes to the `Images` field
- Model must support vision (`qwen2.5vl`, `llava`, `moondream`, etc.)
- Ollama handles encoding internally

Demo opportunity: Show live image description

## Slide: Multimodal: Image Processing

- The `inputs()` function is key
- It splits the prompt on `[img-N]` tags
- Each tag is replaced with embeddings from the vision encoder
- Final sequence interleaves text tokens and image embeddings

This is why the template shows `[img-0]` - it marks where image goes.

## Slide: GGUF: The Model Format

- GGUF = GGML Universal Format (current version 3)
- Previous versions: GGML, GGMF, GGJT

Key features:
- **Single file**: everything in one place
- **Header + metadata + tensors**: clear structure
- **mmap compatible**: load without parsing, OS handles paging
- **Quantization aware**: stores quantization parameters

## Slide: GGUF: Reading Metadata

- The `gguf.Open()` function parses header and metadata
- Useful for inspecting models:
  - Architecture (llama, gemma, qwen, etc.)
  - Context length
  - Vocabulary size
  - Quantization level

Demo: Run `readgguf` tool on a downloaded model

## Slide: GGUF: Tensor Structure

- Tensor naming follows conventions:
  - `token_embd.weight`: vocabulary embeddings
  - `blk.N.*`: transformer block N
  - `output.weight`: final projection

- Understanding this is key for building custom models

## Slide: Build Your Own: Toy Architecture

- This is the advanced section
- We'll create a minimal model from scratch
- It won't produce meaningful text (random weights)
- But it proves you understand the full stack

## Slide: Toy Model: PyTorch Definition

- Simple architecture:
  - Token embeddings
  - N transformer blocks (simplified)
  - Output projection

- Key: the shapes must match what Ollama expects
- Hidden dimension, vocab size, etc. stored in GGUF metadata

## Slide: Toy Model: GGUF Writer

- Python `gguf` package writes the format
- Must set all required metadata:
  - `general.architecture`
  - `*.context_length`
  - `*.embedding_length`
  - `*.block_count`
  - etc.

- Tensor names must match expected patterns
- Can use random weights for testing

## Slide: Toy Model: Creating the Manifest

- `Modelfile` is like Dockerfile for models
- `FROM` specifies base model (or GGUF file)
- Can set parameters, system prompt, template

- `ollama create` registers the model locally

## Slide: Toy Model: Testing

- `ollama run toymodel` should work
- Output will be garbage (untrained weights)
- But the fact it runs proves the format is correct

- Next step: actually train the model or use real weights

## Slide: Complete Request Flow Summary

Recap the entire flow in one diagram:
- Entry: curl/SDK → API
- Middle: Scheduler → Runner subprocess
- Inner loop: tokenize → forward → sample → stream
- Exit: streamed response

## Slide: Key Takeaways

1. Architecture is cleanly separated
2. Templates are critical for model-specific formatting
3. GGUF is the universal container format
4. Multimodal = embeddings interleaved with tokens
5. You can create custom architectures with proper metadata
6. Token sampling has many tunable parameters

## Slide: Resources

- Point to official repos
- Mention the workshop materials for deeper dives
- GGUF spec document is essential reading

## Slide: Questions

- Open for Q&A
- Have the curl examples ready for live demos
- Have a downloaded model ready for `readgguf` demo
