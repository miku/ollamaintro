# Internals

## Ollama Codebase

About 78K LOC Go. Ollama vendors [ggml](https://github.com/ggml-org/ggml) and
[llama.cpp](https://github.com/ggml-org/llama.cpp), which are C and C++ projects.

```
$ tokei
===============================================================================
 Language            Files        Lines         Code     Comments       Blanks
===============================================================================
 GNU Style Assembly      1            6            6            0            0
 Autoconf                1            4            4            0            0
 C                      12        26550        21425          947         4178
 C Header               87       122062       103618        11551         6893
 CMake                   8         1739         1403          112          224
 C++                    58        95684        75391         5559        14734
 C++ Header              3        26105        18285         4684         3136
 CSS                     1           34           29            0            5
 Dockerfile              1          123          104            3           16
 Go                    327        78674        64719         4230         9725
 HTML                    1            9            9            0            0
 JavaScript              2           13           12            1            0
 JSON                   38       147592       147592            0            0
 Objective-C             2         6816         5624          233          959
 PowerShell              2          265          242            9           14
 Protocol Buffers        1          333           98          178           57
 Shell                   7          608          476           39           93
 SVG                     1            9            9            0            0
 Plain Text              4       190212            0       167280        22932
 TSX                     2          129          122            0            7
 TypeScript              9          494          407           14           73
-------------------------------------------------------------------------------
 Markdown               25         3659            0         2439         1220
 |- BASH                 1            1            1            0            0
 |- Dockerfile           3           22           21            0            1
 |- Go                   1           23           22            0            1
 |- INI                  2           18           16            0            2
 |- JavaScript           1           43           34            1            8
 |- JSON                 1          504          504            0            0
 |- PowerShell           3            5            5            0            0
 |- Python               2           96           78            2           16
 |- Shell               15          629          624            0            5
 |- TypeScript           1           18           15            0            3
 (Total)                           5018         1320         2442         1256
===============================================================================
 Total                 593       701120       439575       197279        64266
===============================================================================
```


## Request Flow

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


## GGUF

See also: [GGUF](https://github.com/ggml-org/ggml/blob/master/docs/gguf.md)

> GGUF is a file format for storing models for inference with GGML and
> executors based on GGML. GGUF is a binary format that is designed for fast
> loading and saving of models, and for ease of reading. Models are
> traditionally developed using PyTorch or another framework, and then
> converted to GGUF for use in GGML.

Some format evolution:

> GGML, GGMF, GGJT, and now GGUF — and it is now at **version three** of the format.

### Example: Reading a gguf file

Random, example project:
[SmolDocling-256M-preview-GGUF](https://huggingface.co/Mungert/SmolDocling-256M-preview-GGUF),
[files and
versions](https://huggingface.co/Mungert/SmolDocling-256M-preview-GGUF/tree/main)

```
$ curl -sLo SmolDocling-256M-preview-q8_0.gguf "https://huggingface.co/Mungert/SmolDocling-256M-preview-GGUF/resolve/main/SmolDocling-256M-preview-q8_0.gguf?download=true"
```


