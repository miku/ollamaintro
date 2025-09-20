# Internals

The overall process:

1. **Request Reception**: Client sends HTTP request to Ollama server (port
   11434) via CLI, REST API, or OpenAI-compatible endpoint

2. **Request Routing**: Server routes request through `server/routes.go` to
   appropriate handler (`/api/generate`, `/api/chat`, etc.)

3. **Model Management**: Scheduler (`server/sched.go`) checks if model is
   loaded; if not, loads model into memory with appropriate backend
(CUDA/ROCm/Metal)

4. **Prompt Processing**: Request is parsed and prepared for inference by the
   runner system

5. **Inference Execution**: Specialized runner (`runner/`) performs model
   inference using optimized backends for token generation

6. **Response Streaming**: Generated tokens are streamed back to client via
   HTTP streaming (default mode) or returned as complete response

7. **Resource Management**: System manages GPU/CPU resources and model
   lifecycle for subsequent requests

The entire process runs locally without external API calls, with intelligent
resource scheduling for concurrent model execution.
