# Customization

We are going to customize a model.

* [modelfile docs](https://github.com/ollama/ollama/blob/main/docs/modelfile.md)

## Parameters

* top-p 

> Top-p works by setting a probability threshold p for the next token in an output sequence. The model is allowed to generate responses by using tokens within the probability limit. 

* top-k

> The top-k hyperparameter is another diversity-focused setting. The k value sets the limit for the number of terms that can be considered as the next in the sequence.

* temperature

> A temperature setting of 1 uses the standard probability distribution for the model. Temperatures higher than 1 flatten the probability distribution, encouraging the model to select a wider range of tokens. Conversely, temperatures lower than 1 widen the probability distribution, making the model more likely to select the most probable next token. 

> The difference between temperature and top-p sampling is that while temperature adjusts the probability distribution of potential tokens, top-p sampling limits token selection to a finite group.

Can use top-p and top-k together to limit the number of tokens and then sharpen the distribution.

* system prompt

This prompt get prepended to each request and allows to steer the model towards a interaction mode. Intruction following is  post-trained and models may struggle to actually to what is in the prompt.

## Examples

See [x/custom/](x/custom).

* support for a product
* writing assistant 
* rogerian therapist

You can try it out.

```
$ ollama create NAME -f MODELFILE
```

Example:

``` 
$ ollama create gemma3-eliza-miku-0 -f Modelfile.eliza-miku-0
$ ollama run gemma3-eliza-miku-0 
```


## Prompt engineering

* prompt engineering, context engineering
* more approaches like DSpy or [textgrad](https://arxiv.org/abs/2406.07496)

> Instead of wrangling prompts or training jobs, DSPy (Declarative Self-improving Python) enables you to build AI software from natural-language modules and to generically compose them with different models, inference strategies, or learning algorithms. 

> TEXTGRAD backpropagates textual feedback provided
by LLMs to improve individual components of a compound AI system. 


Sidenote: emergence of tools for assembling prompts, e.g. gitingest, repomix, and many more.

## Evaluation

For a more thorough assessment, you would need to setup an evaluation framework.

Example: [Ollama-MMLU-Pro](https://github.com/chigkim/Ollama-MMLU-Pro) based on [MMLU-Pro](https://huggingface.co/datasets/TIGER-Lab/MMLU-Pro), a variant of [MMLU](https://en.wikipedia.org/wiki/MMLU) "Massive Multitask Language Understanding"

> MMLU consists of 15,908 multiple-choice questions, with 1,540 of them being used to select and assess optimal settings for models – temperature, batch size and learning rate. The questions span across 57 subjects, from highly complex STEM fields and international law, to nutrition and religion. It was one of the most commonly used benchmarks for comparing the capabilities of large language models, with over 100 million downloads as of July 2024.


Example for gemma3:

| overall | biology | business | chemistry | computer science | economics | engineering | health | history | law | math | philosophy | physics | psychology | other |
| ------- | ------- | -------- | --------- | ---------------- | --------- | ----------- | ------ | ------- | --- | ---- | ---------- | ------- | ---------- | ----- |
| 34.23 | 58.58 | 35.23 | 24.03 | 38.29 | 45.14 | 19.81 | 37.04 | 29.13 | 22.07 | 41.45 | 29.26 | 26.40 | 49.50 | 34.31 |


| overall | biology | business | chemistry | computer science | economics | engineering | health | history | law | math | philosophy | physics | psychology | other |
| ------- | ------- | -------- | --------- | ---------------- | --------- | ----------- | ------ | ------- | --- | ---- | ---------- | ------- | ---------- | ----- |
| 46.94 | 58.16 | 55.77 | 52.47 | 55.61 | 56.64 | 33.85 | 39.73 | 24.41 | 19.62 | 69.36 | 35.67 | 51.96 | 51.13 | 35.82 |
