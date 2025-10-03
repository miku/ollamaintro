## Startup

```
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.325+02:00 level=INFO source=routes.go:1475 msg="server config" env="map[CUDA_VISIBLE_DEVICES: GPU_DEVICE_ORDINAL: HIP_VISIBLE_DEVICES: HSA_OVERRIDE_GFX_VERSION: HTTPS_PROXY: HTTP_PROXY: NO_PROXY: OLLAMA_CONTEXT_LENGTH:4096 OLLAMA_DEBUG:DEBUG OLLAMA_FLASH_ATTENTION:false OLLAMA_GPU_OVERHEAD:0 OLLAMA_HOST:http://0.0.0.0:11435 OLLAMA_INTEL_GPU:false OLLAMA_KEEP_ALIVE:5m0s OLLAMA_KV_CACHE_TYPE: OLLAMA_LLM_LIBRARY: OLLAMA_LOAD_TIMEOUT:5m0s OLLAMA_MAX_LOADED_MODELS:0 OLLAMA_MAX_QUEUE:512 OLLAMA_MODELS:/usr/share/ollama/.ollama/models OLLAMA_MULTIUSER_CACHE:false OLLAMA_NEW_ENGINE:false OLLAMA_NOHISTORY:false OLLAMA_NOPRUNE:false OLLAMA_NUM_PARALLEL:1 OLLAMA_ORIGINS:[http://localhost https://localhost http://localhost:* https://localhost:* http://127.0.0.1 https://127.0.0.1 http://127.0.0.1:* https://127.0.0.1:* http://0.0.0.0 https://0.0.0.0 http://0.0.0.0:* https://0.0.0.0:* app://* file://* tauri://* vscode-webview://* vscode-file://*] OLLAMA_REMOTES:[ollama.com] OLLAMA_SCHED_SPREAD:false ROCR_VISIBLE_DEVICES: http_proxy: https_proxy: no_proxy:]"
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.326+02:00 level=INFO source=images.go:518 msg="total blobs: 6"
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.326+02:00 level=INFO source=images.go:525 msg="total unused blobs removed: 0"
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.326+02:00 level=INFO source=routes.go:1528 msg="Listening on [::]:11435 (version 0.12.3)"
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.326+02:00 level=DEBUG source=sched.go:121 msg="starting llm scheduler"
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.326+02:00 level=INFO source=gpu.go:217 msg="looking for compatible GPUs"
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.327+02:00 level=DEBUG source=gpu.go:98 msg="searching for GPU discovery libraries for NVIDIA"
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.327+02:00 level=DEBUG source=gpu.go:520 msg="Searching for GPU library" name=libcuda.so*
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.327+02:00 level=DEBUG source=gpu.go:544 msg="gpu library search" globs="[/usr/lib/ollama/libcuda.so* /libcuda.so* /usr/local/cuda*/targets/*/lib/libcuda.so* /usr/lib/*-linux-gnu/nvidia/current/libcuda.so* /usr/lib/*-linux-gnu/libcuda.so* /usr/lib/wsl/lib/libcuda.so* /usr/lib/wsl/drivers/*/libcuda.so* /opt/cuda/lib*/libcuda.so* /usr/local/cuda/lib*/libcuda.so* /usr/lib*/libcuda.so* /usr/local/lib*/libcuda.so*]"
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.333+02:00 level=DEBUG source=gpu.go:577 msg="discovered GPU libraries" paths=[]
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.333+02:00 level=DEBUG source=gpu.go:520 msg="Searching for GPU library" name=libcudart.so*
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.333+02:00 level=DEBUG source=gpu.go:544 msg="gpu library search" globs="[/usr/lib/ollama/libcudart.so* /libcudart.so* /usr/lib/ollama/cuda_v*/libcudart.so* /usr/local/cuda/lib64/libcudart.so* /usr/lib/x86_64-linux-gnu/nvidia/current/libcudart.so* /usr/lib/x86_64-linux-gnu/libcudart.so* /usr/lib/wsl/lib/libcudart.so* /usr/lib/wsl/drivers/*/libcudart.so* /opt/cuda/lib64/libcudart.so* /usr/local/cuda*/targets/aarch64-linux/lib/libcudart.so* /usr/lib/aarch64-linux-gnu/nvidia/current/libcudart.so* /usr/lib/aarch64-linux-gnu/libcudart.so* /usr/local/cuda/lib*/libcudart.so* /usr/lib*/libcudart.so* /usr/local/lib*/libcudart.so*]"
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.334+02:00 level=DEBUG source=gpu.go:577 msg="discovered GPU libraries" paths="[/usr/lib/ollama/cuda_v12/libcudart.so.12.8.90 /usr/lib/ollama/cuda_v13/libcudart.so.13.0.88]"
Oct 03 13:16:39 gemma ollama[26513]: cudaSetDevice err: 35
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.335+02:00 level=DEBUG source=gpu.go:593 msg="Unable to load cudart library /usr/lib/ollama/cuda_v12/libcudart.so.12.8.90: your nvidia driver is too old or missing.  If you have a CUDA GPU please upgrade to run ollama"
Oct 03 13:16:39 gemma ollama[26513]: cudaSetDevice err: 35
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.335+02:00 level=DEBUG source=gpu.go:593 msg="Unable to load cudart library /usr/lib/ollama/cuda_v13/libcudart.so.13.0.88: your nvidia driver is too old or missing.  If you have a CUDA GPU please upgrade to run ollama"
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.335+02:00 level=DEBUG source=amd_linux.go:423 msg="amdgpu driver not detected /sys/module/amdgpu"
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.335+02:00 level=INFO source=gpu.go:396 msg="no compatible GPUs were discovered"
Oct 03 13:16:39 gemma ollama[26513]: time=2025-10-03T13:16:39.335+02:00 level=INFO source=types.go:131 msg="inference compute" id=0 library=cpu variant="" compute="" driver=0.0 name="" total="31.0 GiB" available="27.4 GiB"
```

## Loading a model

```
Oct 03 13:18:03 gemma ollama[26513]: [GIN] 2025/10/03 - 13:18:03 | 200 |      91.896µs |             ::1 | HEAD     "/"
Oct 03 13:18:03 gemma ollama[26513]: time=2025-10-03T13:18:03.874+02:00 level=DEBUG source=ggml.go:276 msg="key with type not found" key=general.alignment default=32
Oct 03 13:18:03 gemma ollama[26513]: [GIN] 2025/10/03 - 13:18:03 | 200 |    57.95927ms |             ::1 | POST     "/api/show"
Oct 03 13:18:03 gemma ollama[26513]: time=2025-10-03T13:18:03.943+02:00 level=DEBUG source=gpu.go:410 msg="updating system memory data" before.total="31.0 GiB" before.free="27.4 GiB" before.free_swap="977.0 MiB" now.total="31.0 GiB" now.fr
ee="27.5 GiB" now.free_swap="977.0 MiB"
Oct 03 13:18:03 gemma ollama[26513]: time=2025-10-03T13:18:03.943+02:00 level=DEBUG source=sched.go:188 msg="updating default concurrency" OLLAMA_MAX_LOADED_MODELS=3 gpu_count=1
Oct 03 13:18:03 gemma ollama[26513]: time=2025-10-03T13:18:03.954+02:00 level=DEBUG source=ggml.go:276 msg="key with type not found" key=general.alignment default=32
Oct 03 13:18:03 gemma ollama[26513]: time=2025-10-03T13:18:03.954+02:00 level=DEBUG source=sched.go:208 msg="loading first model" model=/usr/share/ollama/.ollama/models/blobs/sha256-dde5aa3fc5ffc17176b5e8bdc82f587b24b2678c6c66101bf7da77af9
f7ccdff
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: loaded meta data with 30 key-value pairs and 255 tensors from /usr/share/ollama/.ollama/models/blobs/sha256-dde5aa3fc5ffc17176b5e8bdc82f587b24b2678c6c66101bf7da77af9f7ccdff (version
GGUF V3 (latest))
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: Dumping metadata keys/values. Note: KV overrides do not apply in this output.
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   0:                       general.architecture str              = llama
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   1:                               general.type str              = model
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   2:                               general.name str              = Llama 3.2 3B Instruct
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   3:                           general.finetune str              = Instruct
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   4:                           general.basename str              = Llama-3.2
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   5:                         general.size_label str              = 3B
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   6:                               general.tags arr[str,6]       = ["facebook", "meta", "pytorch", "llam...
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   7:                          general.languages arr[str,8]       = ["en", "de", "fr", "it", "pt", "hi", ...
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   8:                          llama.block_count u32              = 28
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   9:                       llama.context_length u32              = 131072
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  10:                     llama.embedding_length u32              = 3072
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  11:                  llama.feed_forward_length u32              = 8192
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  12:                 llama.attention.head_count u32              = 24
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  13:              llama.attention.head_count_kv u32              = 8
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  14:                       llama.rope.freq_base f32              = 500000.000000
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  15:     llama.attention.layer_norm_rms_epsilon f32              = 0.000010
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  16:                 llama.attention.key_length u32              = 128
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  17:               llama.attention.value_length u32              = 128
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  18:                          general.file_type u32              = 15
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  19:                           llama.vocab_size u32              = 128256
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  20:                 llama.rope.dimension_count u32              = 128
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  21:                       tokenizer.ggml.model str              = gpt2
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  22:                         tokenizer.ggml.pre str              = llama-bpe
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  23:                      tokenizer.ggml.tokens arr[str,128256]  = ["!", "\"", "#", "$", "%", "&", "'", ...
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  24:                  tokenizer.ggml.token_type arr[i32,128256]  = [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, ...
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  25:                      tokenizer.ggml.merges arr[str,280147]  = ["Ġ Ġ", "Ġ ĠĠĠ", "ĠĠ ĠĠ", "...
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  26:                tokenizer.ggml.bos_token_id u32              = 128000                                                                                                        Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  27:                tokenizer.ggml.eos_token_id u32              = 128009
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  28:                    tokenizer.chat_template str              = {{- bos_token }}\n{%- if custom_tools ...
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  29:               general.quantization_version u32              = 2
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - type  f32:   58 tensors                                                                                                                                                             Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - type q4_K:  168 tensors                                                                                                                                                             Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - type q6_K:   29 tensors
Oct 03 13:18:04 gemma ollama[26513]: print_info: file format = GGUF V3 (latest)                                                                                                                                                                Oct 03 13:18:04 gemma ollama[26513]: print_info: file type   = Q4_K - Medium
Oct 03 13:18:04 gemma ollama[26513]: print_info: file size   = 1.87 GiB (5.01 BPW)
Oct 03 13:18:04 gemma ollama[26513]: init_tokenizer: initializing tokenizer for type 2
Oct 03 13:18:04 gemma ollama[26513]: load: control token: 128254 '<|reserved_special_token_246|>' is not marked as EOG
Oct 03 13:18:04 gemma ollama[26513]: load: control token: 128249 '<|reserved_special_token_241|>' is not marked as EOG
Oct 03 13:18:04 gemma ollama[26513]: load: control token: 128246 '<|reserved_special_token_238|>' is not marked as EOG
Oct 03 13:18:04 gemma ollama[26513]: load: control token: 128243 '<|reserved_special_token_235|>' is not marked as EOG
Oct 03 13:18:04 gemma ollama[26513]: load: control token: 128242 '<|reserved_special_token_234|>' is not marked as EOG
Oct 03 13:18:04 gemma ollama[26513]: load: control token: 128241 '<|reserved_special_token_233|>' is not marked as EOG
Oct 03 13:18:04 gemma ollama[26513]: load: control token: 128240 '<|reserved_special_token_232|>' is not marked as EOG
Oct 03 13:18:04 gemma ollama[26513]: load: control token: 128235 '<|reserved_special_token_227|>' is not marked as EOG
Oct 03 13:18:04 gemma ollama[26513]: load: control token: 128231 '<|reserved_special_token_223|>' is not marked as EOG
Oct 03 13:18:04 gemma ollama[26513]: load: control token: 128230 '<|reserved_special_token_222|>' is not marked as EOG

...

Oct 03 13:18:04 gemma ollama[26513]: load: printing all EOG tokens:
Oct 03 13:18:04 gemma ollama[26513]: load:   - 128001 ('<|end_of_text|>')
Oct 03 13:18:04 gemma ollama[26513]: load:   - 128008 ('<|eom_id|>')
Oct 03 13:18:04 gemma ollama[26513]: load:   - 128009 ('<|eot_id|>')
Oct 03 13:18:04 gemma ollama[26513]: load: special tokens cache size = 256
Oct 03 13:18:04 gemma ollama[26513]: load: token to piece cache size = 0.7999 MB
Oct 03 13:18:04 gemma ollama[26513]: print_info: arch             = llama
Oct 03 13:18:04 gemma ollama[26513]: print_info: vocab_only       = 1
Oct 03 13:18:04 gemma ollama[26513]: print_info: model type       = ?B
Oct 03 13:18:04 gemma ollama[26513]: print_info: model params     = 3.21 B
Oct 03 13:18:04 gemma ollama[26513]: print_info: general.name     = Llama 3.2 3B Instruct
Oct 03 13:18:04 gemma ollama[26513]: print_info: vocab type       = BPE
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_vocab          = 128256
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_merges         = 280147
Oct 03 13:18:04 gemma ollama[26513]: print_info: BOS token        = 128000 '<|begin_of_text|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: EOS token        = 128009 '<|eot_id|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: EOT token        = 128009 '<|eot_id|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: EOM token        = 128008 '<|eom_id|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: LF token         = 198 'Ċ'
Oct 03 13:18:04 gemma ollama[26513]: print_info: EOG token        = 128001 '<|end_of_text|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: EOG token        = 128008 '<|eom_id|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: EOG token        = 128009 '<|eot_id|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: max token length = 256
Oct 03 13:18:04 gemma ollama[26513]: llama_model_load: vocab only - skipping tensors
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.222+02:00 level=DEBUG source=gpu.go:410 msg="updating system memory data" before.total="31.0 GiB" before.free="27.5 GiB" before.free_swap="977.0 MiB" now.total="31.0 GiB" now.fr
ee="27.4 GiB" now.free_swap="977.0 MiB"
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.223+02:00 level=INFO source=server.go:399 msg="starting runner" cmd="/usr/bin/ollama runner --model /usr/share/ollama/.ollama/models/blobs/sha256-dde5aa3fc5ffc17176b5e8bdc82f587
b24b2678c6c66101bf7da77af9f7ccdff --port 36491"
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.223+02:00 level=DEBUG source=server.go:400 msg=subprocess PATH=$PATH OLLAMA_HOST=0.0.0.0:11435 OLLAMA_DEBUG=1 OLLAMA_MAX_LOADED_MODELS=3 OLLAMA_LIBRARY_PATH=/usr/lib/ollama LD_L
IBRARY_PATH=/usr/lib/ollama:/usr/lib/ollama
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.223+02:00 level=DEBUG source=gpu.go:410 msg="updating system memory data" before.total="31.0 GiB" before.free="27.4 GiB" before.free_swap="977.0 MiB" now.total="31.0 GiB" now.fr
ee="27.4 GiB" now.free_swap="977.0 MiB"
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.223+02:00 level=INFO source=server.go:504 msg="system memory" total="31.0 GiB" free="27.4 GiB" free_swap="977.0 MiB"
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.223+02:00 level=DEBUG source=memory.go:181 msg=evaluating library=cpu gpu_count=1 available="[27.5 GiB]"
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.223+02:00 level=DEBUG source=ggml.go:276 msg="key with type not found" key=llama.vision.block_count default=0


Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.223+02:00 level=DEBUG source=ggml.go:611 msg="default cache size estimate" "attention MiB"=448 "attention bytes"=469762048 "recurrent MiB"=0 "recurrent bytes"=0
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.223+02:00 level=INFO source=memory.go:36 msg="new model will fit in available VRAM across minimum required GPUs, loading" model=/usr/share/ollama/.ollama/models/blobs/sha256-dde
5aa3fc5ffc17176b5e8bdc82f587b24b2678c6c66101bf7da77af9f7ccdff library=cpu parallel=1 required="0 B" gpus=1
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.223+02:00 level=DEBUG source=memory.go:181 msg=evaluating library=cpu gpu_count=1 available="[27.5 GiB]"
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.223+02:00 level=DEBUG source=ggml.go:276 msg="key with type not found" key=llama.vision.block_count default=0
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.224+02:00 level=DEBUG source=ggml.go:611 msg="default cache size estimate" "attention MiB"=448 "attention bytes"=469762048 "recurrent MiB"=0 "recurrent bytes"=0
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.224+02:00 level=INFO source=server.go:544 msg=offload library=cpu layers.requested=-1 layers.model=29 layers.offload=0 layers.split=[] memory.available="[27.5 GiB]" memory.gpu_o
verhead="0 B" memory.required.full="2.6 GiB" memory.required.partial="0 B" memory.required.kv="448.0 MiB" memory.required.allocations="[2.6 GiB]" memory.weights.total="1.9 GiB" memory.weights.repeating="1.6 GiB" memory.weights.nonrepeating
="308.2 MiB" memory.graph.full="256.5 MiB" memory.graph.partial="570.7 MiB"
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.232+02:00 level=INFO source=runner.go:864 msg="starting go runner"
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.232+02:00 level=DEBUG source=ggml.go:94 msg="ggml backend load all from path" path=/usr/lib/ollama
Oct 03 13:18:04 gemma ollama[26513]: load_backend: loaded CPU backend from /usr/lib/ollama/libggml-cpu-alderlake.so
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.236+02:00 level=INFO source=ggml.go:104 msg=system CPU.0.SSE3=1 CPU.0.SSSE3=1 CPU.0.AVX=1 CPU.0.AVX_VNNI=1 CPU.0.AVX2=1 CPU.0.F16C=1 CPU.0.FMA=1 CPU.0.BMI2=1 CPU.0.LLAMAFILE=1 C
PU.1.LLAMAFILE=1 compiler=cgo(gcc)
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.236+02:00 level=INFO source=runner.go:900 msg="Server listening on 127.0.0.1:36491"
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.246+02:00 level=INFO source=runner.go:799 msg=load request="{Operation:commit LoraPath:[] Parallel:1 BatchSize:512 FlashAttention:false KvSize:4096 KvCacheType: NumThreads:2 GPU
Layers:[] MultiUserCache:false ProjectorPath: MainGPU:0 UseMmap:false}"
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.246+02:00 level=INFO source=server.go:1251 msg="waiting for llama runner to start responding"
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.247+02:00 level=INFO source=server.go:1285 msg="waiting for server to become available" status="llm server loading model"
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: loaded meta data with 30 key-value pairs and 255 tensors from /usr/share/ollama/.ollama/models/blobs/sha256-dde5aa3fc5ffc17176b5e8bdc82f587b24b2678c6c66101bf7da77af9f7ccdff (version
GGUF V3 (latest))
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: Dumping metadata keys/values. Note: KV overrides do not apply in this output.
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   0:                       general.architecture str              = llama
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   1:                               general.type str              = model
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   2:                               general.name str              = Llama 3.2 3B Instruct
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   3:                           general.finetune str              = Instruct
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   4:                           general.basename str              = Llama-3.2
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   5:                         general.size_label str              = 3B
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   6:                               general.tags arr[str,6]       = ["facebook", "meta", "pytorch", "llam...
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   7:                          general.languages arr[str,8]       = ["en", "de", "fr", "it", "pt", "hi", ...
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   8:                          llama.block_count u32              = 28
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv   9:                       llama.context_length u32              = 131072
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  10:                     llama.embedding_length u32              = 3072
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  11:                  llama.feed_forward_length u32              = 8192
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  12:                 llama.attention.head_count u32              = 24
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  13:              llama.attention.head_count_kv u32              = 8
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  14:                       llama.rope.freq_base f32              = 500000.000000
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  15:     llama.attention.layer_norm_rms_epsilon f32              = 0.000010
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  16:                 llama.attention.key_length u32              = 128
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  17:               llama.attention.value_length u32              = 128
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  18:                          general.file_type u32              = 15
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  19:                           llama.vocab_size u32              = 128256
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  20:                 llama.rope.dimension_count u32              = 128
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  21:                       tokenizer.ggml.model str              = gpt2
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  22:                         tokenizer.ggml.pre str              = llama-bpe
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  23:                      tokenizer.ggml.tokens arr[str,128256]  = ["!", "\"", "#", "$", "%", "&", "'", ...
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  24:                  tokenizer.ggml.token_type arr[i32,128256]  = [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, ...
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  25:                      tokenizer.ggml.merges arr[str,280147]  = ["Ġ Ġ", "Ġ ĠĠĠ", "ĠĠ ĠĠ", "...
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  26:                tokenizer.ggml.bos_token_id u32              = 128000
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  27:                tokenizer.ggml.eos_token_id u32              = 128009
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  28:                    tokenizer.chat_template str              = {{- bos_token }}\n{%- if custom_tools ...
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - kv  29:               general.quantization_version u32              = 2

...

Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - type  f32:   58 tensors
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - type q4_K:  168 tensors
Oct 03 13:18:04 gemma ollama[26513]: llama_model_loader: - type q6_K:   29 tensors
Oct 03 13:18:04 gemma ollama[26513]: print_info: file format = GGUF V3 (latest)
Oct 03 13:18:04 gemma ollama[26513]: print_info: file type   = Q4_K - Medium
Oct 03 13:18:04 gemma ollama[26513]: print_info: file size   = 1.87 GiB (5.01 BPW)
Oct 03 13:18:04 gemma ollama[26513]: init_tokenizer: initializing tokenizer for type 2

...

Oct 03 13:18:04 gemma ollama[26513]: load: printing all EOG tokens:
Oct 03 13:18:04 gemma ollama[26513]: load:   - 128001 ('<|end_of_text|>')
Oct 03 13:18:04 gemma ollama[26513]: load:   - 128008 ('<|eom_id|>')
Oct 03 13:18:04 gemma ollama[26513]: load:   - 128009 ('<|eot_id|>')
Oct 03 13:18:04 gemma ollama[26513]: load: special tokens cache size = 256
Oct 03 13:18:04 gemma ollama[26513]: load: token to piece cache size = 0.7999 MB
Oct 03 13:18:04 gemma ollama[26513]: print_info: arch             = llama
Oct 03 13:18:04 gemma ollama[26513]: print_info: vocab_only       = 0
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_ctx_train      = 131072
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_embd           = 3072
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_layer          = 28
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_head           = 24
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_head_kv        = 8
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_rot            = 128
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_swa            = 0
Oct 03 13:18:04 gemma ollama[26513]: print_info: is_swa_any       = 0
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_embd_head_k    = 128
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_embd_head_v    = 128
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_gqa            = 3
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_embd_k_gqa     = 1024
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_embd_v_gqa     = 1024
Oct 03 13:18:04 gemma ollama[26513]: print_info: f_norm_eps       = 0.0e+00
Oct 03 13:18:04 gemma ollama[26513]: print_info: f_norm_rms_eps   = 1.0e-05
Oct 03 13:18:04 gemma ollama[26513]: print_info: f_clamp_kqv      = 0.0e+00
Oct 03 13:18:04 gemma ollama[26513]: print_info: f_max_alibi_bias = 0.0e+00
Oct 03 13:18:04 gemma ollama[26513]: print_info: f_logit_scale    = 0.0e+00
Oct 03 13:18:04 gemma ollama[26513]: print_info: f_attn_scale     = 0.0e+00
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_ff             = 8192
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_expert         = 0
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_expert_used    = 0
Oct 03 13:18:04 gemma ollama[26513]: print_info: causal attn      = 1
Oct 03 13:18:04 gemma ollama[26513]: print_info: pooling type     = 0
Oct 03 13:18:04 gemma ollama[26513]: print_info: rope type        = 0
Oct 03 13:18:04 gemma ollama[26513]: print_info: rope scaling     = linear
Oct 03 13:18:04 gemma ollama[26513]: print_info: freq_base_train  = 500000.0
Oct 03 13:18:04 gemma ollama[26513]: print_info: freq_scale_train = 1
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_ctx_orig_yarn  = 131072
Oct 03 13:18:04 gemma ollama[26513]: print_info: rope_finetuned   = unknown
Oct 03 13:18:04 gemma ollama[26513]: print_info: model type       = 3B
Oct 03 13:18:04 gemma ollama[26513]: print_info: model params     = 3.21 B
Oct 03 13:18:04 gemma ollama[26513]: print_info: general.name     = Llama 3.2 3B Instruct
Oct 03 13:18:04 gemma ollama[26513]: print_info: vocab type       = BPE
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_vocab          = 128256
Oct 03 13:18:04 gemma ollama[26513]: print_info: n_merges         = 280147
Oct 03 13:18:04 gemma ollama[26513]: print_info: BOS token        = 128000 '<|begin_of_text|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: EOS token        = 128009 '<|eot_id|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: EOT token        = 128009 '<|eot_id|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: EOM token        = 128008 '<|eom_id|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: LF token         = 198 'Ċ'
Oct 03 13:18:04 gemma ollama[26513]: print_info: EOG token        = 128001 '<|end_of_text|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: EOG token        = 128008 '<|eom_id|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: EOG token        = 128009 '<|eot_id|>'
Oct 03 13:18:04 gemma ollama[26513]: print_info: max token length = 256
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: loading model tensors, this can take a while... (mmap = false)

...

Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer   0 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer   1 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer   2 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer   3 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer   4 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer   5 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer   6 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer   7 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer   8 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer   9 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  10 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  11 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  12 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  13 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  14 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  15 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  16 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  17 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  18 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  19 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  20 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  21 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  22 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  23 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  24 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  25 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  26 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  27 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors: layer  28 assigned to device CPU, is_swa = 0
Oct 03 13:18:04 gemma ollama[26513]: load_tensors:          CPU model buffer size =  1918.35 MiB
Oct 03 13:18:04 gemma ollama[26513]: load_all_data: no device found for buffer type CPU for async uploads
Oct 03 13:18:04 gemma ollama[26513]: time=2025-10-03T13:18:04.749+02:00 level=DEBUG source=server.go:1295 msg="model load progress 0.70"
Oct 03 13:18:04 gemma ollama[26513]: llama_context: constructing llama_context
Oct 03 13:18:04 gemma ollama[26513]: llama_context: n_seq_max     = 1
Oct 03 13:18:04 gemma ollama[26513]: llama_context: n_ctx         = 4096
Oct 03 13:18:04 gemma ollama[26513]: llama_context: n_ctx_per_seq = 4096
Oct 03 13:18:04 gemma ollama[26513]: llama_context: n_batch       = 512
Oct 03 13:18:04 gemma ollama[26513]: llama_context: n_ubatch      = 512
Oct 03 13:18:04 gemma ollama[26513]: llama_context: causal_attn   = 1
Oct 03 13:18:04 gemma ollama[26513]: llama_context: flash_attn    = 0
Oct 03 13:18:04 gemma ollama[26513]: llama_context: kv_unified    = false
Oct 03 13:18:04 gemma ollama[26513]: llama_context: freq_base     = 500000.0
Oct 03 13:18:04 gemma ollama[26513]: llama_context: freq_scale    = 1

...

Oct 03 13:18:04 gemma ollama[26513]: llama_context: n_ctx_per_seq (4096) < n_ctx_train (131072) -- the full capacity of the model will not be utilized
Oct 03 13:18:04 gemma ollama[26513]: set_abort_callback: call
Oct 03 13:18:04 gemma ollama[26513]: llama_context:        CPU  output buffer size =     0.50 MiB
Oct 03 13:18:04 gemma ollama[26513]: create_memory: n_ctx = 4096 (padded)
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer   0: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer   1: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer   2: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer   3: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer   4: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer   5: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer   6: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer   7: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer   8: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer   9: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  10: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  11: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  12: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  13: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  14: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  15: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  16: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  17: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  18: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  19: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  20: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  21: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  22: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  23: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  24: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  25: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  26: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: layer  27: dev = CPU
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified:        CPU KV buffer size =   448.00 MiB
Oct 03 13:18:04 gemma ollama[26513]: llama_kv_cache_unified: size =  448.00 MiB (  4096 cells,  28 layers,  1/1 seqs), K (f16):  224.00 MiB, V (f16):  224.00 MiB
Oct 03 13:18:04 gemma ollama[26513]: llama_context: enumerating backends
Oct 03 13:18:04 gemma ollama[26513]: llama_context: backend_ptrs.size() = 1
Oct 03 13:18:04 gemma ollama[26513]: llama_context: max_nodes = 2040
Oct 03 13:18:04 gemma ollama[26513]: llama_context: worst-case: n_tokens = 512, n_seqs = 1, n_outputs = 0
Oct 03 13:18:04 gemma ollama[26513]: graph_reserve: reserving a graph for ubatch with n_tokens =  512, n_seqs =  1, n_outputs =  512
Oct 03 13:18:04 gemma ollama[26513]: graph_reserve: reserving a graph for ubatch with n_tokens =    1, n_seqs =  1, n_outputs =    1
Oct 03 13:18:04 gemma ollama[26513]: graph_reserve: reserving a graph for ubatch with n_tokens =  512, n_seqs =  1, n_outputs =  512
Oct 03 13:18:04 gemma ollama[26513]: llama_context:        CPU compute buffer size =   256.50 MiB
Oct 03 13:18:04 gemma ollama[26513]: llama_context: graph nodes  = 986
Oct 03 13:18:04 gemma ollama[26513]: llama_context: graph splits = 1

...

Oct 03 13:18:05 gemma ollama[26513]: time=2025-10-03T13:18:05.001+02:00 level=INFO source=server.go:1289 msg="llama runner started in 0.78 seconds"
Oct 03 13:18:05 gemma ollama[26513]: time=2025-10-03T13:18:05.001+02:00 level=INFO source=sched.go:470 msg="loaded runners" count=1
Oct 03 13:18:05 gemma ollama[26513]: time=2025-10-03T13:18:05.001+02:00 level=INFO source=server.go:1251 msg="waiting for llama runner to start responding"
Oct 03 13:18:05 gemma ollama[26513]: time=2025-10-03T13:18:05.002+02:00 level=INFO source=server.go:1289 msg="llama runner started in 0.78 seconds"
Oct 03 13:18:05 gemma ollama[26513]: time=2025-10-03T13:18:05.002+02:00 level=DEBUG source=sched.go:482 msg="finished setting up" runner.name=registry.ollama.ai/library/llama3.2:latest runner.inference=cpu runner.devices=1 runner.size="2.6
 GiB" runner.vram="0 B" runner.parallel=1 runner.pid=27080 runner.model=/usr/share/ollama/.ollama/models/blobs/sha256-dde5aa3fc5ffc17176b5e8bdc82f587b24b2678c6c66101bf7da77af9f7ccdff runner.num_ctx=4096
Oct 03 13:18:05 gemma ollama[26513]: [GIN] 2025/10/03 - 13:18:05 | 200 |  1.125268431s |             ::1 | POST     "/api/generate"
Oct 03 13:18:05 gemma ollama[26513]: time=2025-10-03T13:18:05.002+02:00 level=DEBUG source=sched.go:490 msg="context for request finished"
Oct 03 13:18:05 gemma ollama[26513]: time=2025-10-03T13:18:05.002+02:00 level=DEBUG source=sched.go:286 msg="runner with non-zero duration has gone idle, adding timer" runner.name=registry.ollama.ai/library/llama3.2:latest runner.inference
=cpu runner.devices=1 runner.size="2.6 GiB" runner.vram="0 B" runner.parallel=1 runner.pid=27080 runner.model=/usr/share/ollama/.ollama/models/blobs/sha256-dde5aa3fc5ffc17176b5e8bdc82f587b24b2678c6c66101bf7da77af9f7ccdff runner.num_ctx=409
6 duration=5m0s
Oct 03 13:18:05 gemma ollama[26513]: time=2025-10-03T13:18:05.003+02:00 level=DEBUG source=sched.go:304 msg="after processing request finished event" runner.name=registry.ollama.ai/library/llama3.2:latest runner.inference=cpu runner.device
s=1 runner.size="2.6 GiB" runner.vram="0 B" runner.parallel=1 runner.pid=27080 runner.model=/usr/share/ollama/.ollama/models/blobs/sha256-dde5aa3fc5ffc17176b5e8bdc82f587b24b2678c6c66101bf7da77af9f7ccdff runner.num_ctx=4096 refCount=0
```

## Running a completion

```
Oct 03 13:22:29 gemma ollama[26513]: [GIN] 2025/10/03 - 13:22:29 | 200 |      35.456µs |             ::1 | HEAD     "/"
Oct 03 13:22:29 gemma ollama[26513]: time=2025-10-03T13:22:29.933+02:00 level=DEBUG source=ggml.go:276 msg="key with type not found" key=general.alignment default=32
Oct 03 13:22:29 gemma ollama[26513]: [GIN] 2025/10/03 - 13:22:29 | 200 |   58.228479ms |             ::1 | POST     "/api/show"
Oct 03 13:22:30 gemma ollama[26513]: time=2025-10-03T13:22:30.006+02:00 level=DEBUG source=sched.go:580 msg="evaluating already loaded" model=/usr/share/ollama/.ollama/models/blobs/sha256-dde5aa3fc5ffc17176b5e8bdc82f587b24b2678c6c66101bf7da77af9f7ccdff
Oct 03 13:22:30 gemma ollama[26513]: time=2025-10-03T13:22:30.007+02:00 level=DEBUG source=server.go:1388 msg="completion request" images=0 prompt=197 format=""
Oct 03 13:22:30 gemma ollama[26513]: time=2025-10-03T13:22:30.008+02:00 level=DEBUG source=cache.go:104 msg="loading cache slot" id=0 cache=0 prompt=26 used=0 remaining=26
Oct 03 13:22:32 gemma ollama[26513]: [GIN] 2025/10/03 - 13:22:32 | 200 |   2.10313825s |             ::1 | POST     "/api/generate"
Oct 03 13:22:32 gemma ollama[26513]: time=2025-10-03T13:22:32.038+02:00 level=DEBUG source=sched.go:377 msg="context for request finished" runner.name=registry.ollama.ai/library/llama3.2:latest runner.inference=cpu runner.devices=1 runner.size="2.6 GiB" runner.vram="0 B" runner.parallel=1 runner.pid=27080 runner.model=/usr/share/ollama/.ollama/models/blobs/sha256-dde5aa3fc5ffc17176b5e8bdc82f587b24b2678c6c66101bf7da77af9f7ccdff runner.num_ctx=4096
Oct 03 13:22:32 gemma ollama[26513]: time=2025-10-03T13:22:32.038+02:00 level=DEBUG source=sched.go:286 msg="runner with non-zero duration has gone idle, adding timer" runner.name=registry.ollama.ai/library/llama3.2:latest runner.inference=cpu runner.devices=1 runner.size="2.6 GiB" runner.vram="0 B" runner.parallel=1 runner.pid=27080 runner.model=/usr/share/ollama/.ollama/models/blobs/sha256-dde5aa3fc5ffc17176b5e8bdc82f587b24b2678c6c66101bf7da77af9f7ccdff runner.num_ctx=4096 duration=5m0s
Oct 03 13:22:32 gemma ollama[26513]: time=2025-10-03T13:22:32.038+02:00 level=DEBUG source=sched.go:304 msg="after processing request finished event" runner.name=registry.ollama.ai/library/llama3.2:latest runner.inference=cpu runner.devices=1 runner.size="2.6 GiB" runner.vram="0 B" runner.parallel=1 runner.pid=27080 runner.model=/usr/share/ollama/.ollama/models/blobs/sha256-dde5aa3fc5ffc17176b5e8bdc82f587b24b2678c6c66101bf7da77af9f7ccdff runner.num_ctx=4096 refCount=0
```
