# Model Types

Ollama categories as of 09/2025:

* cloud (hosted)
* embedding
* vision
* tools
* thinking

## Embedding

* maps data (text, image, ...) into a dense vector space
* similar documents are mapped closer together

----

Example: embeddinggemma

> EmbeddingGemma is a 308M parameter multilingual text embedding model based on
> Gemma 3. It is optimized for use in everyday devices, such as phones,
> laptops, and tablets. The model produces numerical representations of text to
> be used for downstream tasks like information retrieval, semantic similarity
> search, classification, and clustering. (cf. [model card](https://ai.google.dev/gemma/docs/embeddinggemma/model_card))

Input:

* Text string, such as a question, a prompt, or a document to be embedded (can use [prompt instructions](https://ai.google.dev/gemma/docs/embeddinggemma/model_card#prompt-instructions))
* Maximum input context length of 2K

Output:

* Numerical vector representations of input text data
* Output embedding dimension size of 768, with smaller options available (512, 256, or 128) -- [EmbeddingGemma: Powerful and Lightweight Text Representations](https://arxiv.org/abs/2509.20354)

----

* embedding models are computationally less expensive (can run on low power hardware, w/o GPU)
* one benchmark for embedding models is MTEB ("Massive Text Embedding
  Benchmark",
[https://arxiv.org/pdf/2210.07316](https://arxiv.org/pdf/2210.07316), "MTEB
spans 8 embedding tasks covering a total of 58 datasets and 112 languages", 8
categories: "Bitext mining, classi- fication, clustering, pair classification,
reranking, retrieval, STS and summarization", [leaderboard](https://huggingface.co/spaces/mteb/leaderboard))
* example:

## Vision

*

## Embedding models

* smaller models, with a single vector output; usually faster and can run on
  slower hardware as well

## Image understanding

* llava
* qwenvl

## Reasoning

* additional expansion of the user query intro a search tree
