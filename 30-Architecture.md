# Architecture

Layered architecture.

* API layer
* Scheduler and orchestrator
* Runner layer

## API layer

* HTTP server (port 11434)
* REST endpoints + OpenAI compatibility
* Request routing and validation

## Scheduler

* Model lifecycle management
* Request queuing and prioritization
* Resource allocation (GPU/CPU/Memory)
* Hot-swapping models when memory constrained

## Runner layer

* Subprocess per model
* Two implementations: llamarunner (llama.cpp) and ollamarunner (direct ggml)
* Handles actual inference

## Characteristics

* Docker-inspired UX (pull, run, push)
* Content-addressable storage (like git, docker, ...)
* Streaming responses by default
* Automatic resource management
