# What Happens When You Type a Prompt into Ollama?

> A 20-30 minute technical presentation condensed from a 3-hour workshop

## Overview

This presentation traces the complete journey of a prompt through Ollama's internals, from API request to streamed response. It's designed for a technical audience familiar with basic ML concepts.

## Contents

```
2026-02-18-kimhs/
├── README.md           # This file
├── slides.md           # Presentation slides (markdown)
├── notes.md            # Speaker notes
└── examples/
    ├── 01_completion.go       # Basic text completion
    ├── 02_multimodal.go       # Vision model example
    ├── 03_embedding.go        # Embedding generation
    ├── 04_gguf_reader.go      # Reading GGUF files
    ├── 05_toy_model.py        # Create toy architecture with random weights
    ├── 06_debug_template.sh   # Template debugging demo
    └── Modelfile.example      # Custom model configuration
```

## Key Topics Covered

1. **Architecture Overview** - Three-layer design: API → Scheduler → Runner
2. **Request Journey** - 8 steps from prompt to response
3. **Multimodal** - How vision models process images
4. **GGUF Format** - The universal model container
5. **Build Your Own** - Creating a toy architecture with GGML

## Running the Examples

### Prerequisites

```bash
# Ensure Ollama is running
ollama serve

# Pull required models
ollama pull gemma3:270m
ollama pull embeddinggemma
ollama pull qwen2.5vl  # for multimodal examples
```

### Go Examples

```bash
cd examples

# Basic completion
go run 01_completion.go

# Multimodal (needs an image file)
go run 02_multimodal.go path/to/image.png

# Embeddings
go run 03_embedding.go

# GGUF reader (needs a model file)
go run 04_gguf_reader.go ~/.ollama/models/blobs/sha256-<hash>
```

### Python Examples

```bash
# Install gguf package
pip install gguf numpy

# Create toy model
python 05_toy_model.py

# Import into Ollama
ollama create toymodel -f toymodel.gguf
ollama run toymodel "Hello"
```

### Shell Scripts

```bash
# Debug template rendering
chmod +x 06_debug_template.sh
./06_debug_template.sh gemma3:270m
```

## Presentation Tips

1. **Start with the curl demo** - Show a live request to make it concrete
2. **Use `_debug_render_only`** - Demystifies template processing
3. **Run `readgguf`** - Shows the model format is inspectable
4. **Create the toy model live** - Proves understanding of the full stack

## Time Allocation (25 min)

| Section | Time |
|---------|------|
| Architecture Overview | 3 min |
| Request Journey | 8 min |
| Multimodal | 5 min |
| GGUF Format | 4 min |
| Build Your Own | 5 min |

## Resources

- [Ollama GitHub](https://github.com/ollama/ollama)
- [GGUF Specification](https://github.com/ggml-org/ggml/blob/master/docs/gguf.md)
- [llama.cpp](https://github.com/ggml-org/llama.cpp)
- [Full Workshop](https://github.com/miku/ollamaintro)
