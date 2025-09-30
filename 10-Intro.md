# Intro

## Quick poll

* [ ] Who has already used Ollama?
* [ ] Who has run an LLM locally?
* [ ] Who has worked with model APIs before?

## Workshop goals

* [ ] Better understanding how Ollama works under the hood
* [ ] Learn about model lifecycle and customization
* [ ] Build practical applications with embeddings in Go
* [ ] Create a visual search tool for handwritten notes

## Local and remote machine setup

* Install ollama on your laptop, you can follow the instructions here: [https://ollama.com/download](https://ollama.com/download)
* We will have access to a cloud GPU server instance, running a RTX 4000 SFF GPU with 20GB VRAM (an efficient GPU, that consumes at most 70W)

Most examples will work well on Linux and MacOS, but Windows WSL should also
work.

## Checkpoint: installation

```
$ ollama --version
$ ollama ls
$ curl localhost:11434
```

## Workshop material

All material available on
[github.com/miku/ollamaintro](https://github.com/miku/ollamaintro). Please
clone, as if contains also code examples and project scaffolding.

The workshop was developed on various Linux systems and tested on low end CPU
and a higher end GPU equipped systems:

* [Intel N150 CPU](https://www.intel.de/content/www/de/de/products/sku/241636/intel-processor-n150-6m-cache-up-to-3-60-ghz/specifications.html), 6W TDP
* [Intel i5-1345U](https://www.intel.com/content/www/us/en/products/sku/232127/intel-core-i51345u-processor-12m-cache-up-to-4-70-ghz/specifications.html), 15W TDP
* [Intel i9-13900T](https://www.intel.com/content/www/us/en/products/sku/230498/intel-core-i913900t-processor-36m-cache-up-to-5-30-ghz/compatible.html), 35W TDP, plus [NVIDIA RTX 4000 SFF](https://www.nvidia.com/en-us/products/workstations/rtx-4000-sff/) (70W)

TODO: [AMD Ryzen AI Max+ 395](https://www.amd.com/en/products/processors/laptop/ryzen/ai-300-series/amd-ryzen-ai-max-plus-395.html)

## ollama community integrations

* [community integrations](https://github.com/ollama/ollama?tab=readme-ov-file#community-integrations)

## ollama installation files (linux)

```
$ tar -tf ollama-linux-amd64.tgz
./
./bin/
./bin/ollama
./lib/
./lib/ollama/
./lib/ollama/libggml-hip.so
./lib/ollama/libggml-cpu-sse42.so
./lib/ollama/libggml-base.so
./lib/ollama/libggml-cpu-skylakex.so
./lib/ollama/libggml-cpu-x64.so
./lib/ollama/libggml-cpu-icelake.so
./lib/ollama/cuda_v12/
./lib/ollama/cuda_v12/libcublasLt.so.12
./lib/ollama/cuda_v12/libcublasLt.so.12.8.4.1
./lib/ollama/cuda_v12/libggml-cuda.so
./lib/ollama/cuda_v12/libcudart.so.12
./lib/ollama/cuda_v12/libcudart.so.12.8.90
./lib/ollama/cuda_v12/libcublas.so.12.8.4.1
./lib/ollama/cuda_v12/libcublas.so.12
./lib/ollama/libggml-cpu-sandybridge.so
./lib/ollama/libggml-cpu-alderlake.so
./lib/ollama/cuda_v13/
./lib/ollama/cuda_v13/libcublasLt.so.13
./lib/ollama/cuda_v13/libcudart.so.13.0.88
./lib/ollama/cuda_v13/libggml-cuda.so
./lib/ollama/cuda_v13/libcublas.so.13.0.2.14
./lib/ollama/cuda_v13/libcudart.so.13
./lib/ollama/cuda_v13/libcublasLt.so.13.0.2.14
./lib/ollama/cuda_v13/libcublas.so.13
./lib/ollama/libggml-cpu-haswell.so
```

* libcurart: [cuda runtime docs](https://docs.nvidia.com/cuda/cuda-c-programming-guide/#cuda-runtime)
* libcublast: "The cuBLASLt is a lightweight library dedicated to GEneral
  Matrix-to-matrix Multiply (GEMM) operations with a new flexible API." -- [docs](https://docs.nvidia.com/cuda/cublas/index.html#using-the-cublasLt-api)
* libcublas: [](https://developer.nvidia.com/cublas), "[...] several API
  extensions for providing drop-in industry standard BLAS APIs and GEMM APIs
with support for fusions that are highly optimized for NVIDIA GPUs."
