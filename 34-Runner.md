# Runner

Ollama serve subcommand starts an http server that coordinates requests and
spawns a subprocess per model to serve the actual requests.

This is the serve process and one model loaded from the blob store.

```
$ pstree -apT 143757
ollama,143757 serve
  └─ollama,2720443 runner --ollama-engine --model /usr/share/ollama/.ollama/models/blobs/sha256-3d0b790534fe4b79525fc3692950408dca41171676ed7e21db57af5c65ef6ab6 --port 36431
```

We can run multiple models side by side, until we run out of RAM.

```
$ pstree -apT 143757
ollama,143757 serve
  ├─ollama,2720443 runner --ollama-engine --model /usr/share/ollama/.ollama/models/blobs/sha256-3d0b790534fe4b79525fc3692950408dca41171676ed7e21db57af5c65ef6ab6 --port 36431
  ├─ollama,3088461 runner --model /usr/share/ollama/.ollama/models/blobs/sha256-74701a8c35f6c8d9a4b91f3f3497643001d63e0c7a84e085bed452548fa88d45 --port 38555
  └─ollama,3089161 runner --ollama-engine --model /usr/share/ollama/.ollama/models/blobs/sha256-aeda25e63ebd698fab8638ffb778e68bed908b960d39d0becc650fa981609d25 --port 43731
```

The runner process is a slim server:

```
$ /usr/bin/ollama runner --help
Runner usage
  -model string
        Path to model binary file
  -port int
        Port to expose the server on (default 8080)
  -verbose
        verbose output (default: disabled)
```

Example invocation:

```
$ /usr/bin/ollama runner --ollama-engine --model /usr/share/ollama/.ollama/models/blobs/sha256-3d0b790534fe4b79525fc3692950408dca41171676ed7e21db57af5c65ef6ab6 --port 36432
```

## Background

Ollama started with a wrapper around llama and since 02/2025 (cf.
[#7913](https://github.com/ollama/ollama/pull/7913)) comes with an additional
runner type "ollama", that interfaces with ggml directly.

Currently two runner types:

* classic (cgo, llama.cpp)
* ollama (cgo, ggml)



