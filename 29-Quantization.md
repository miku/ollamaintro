## Quantization

* model size depends on number of parameters, typically millions, billions, up to trillions
* one technique to reduce model size (other: lora, pruning, knowledge
  destillation, parameter sharing; cf. [model shrinking
techniques](https://web.dev/articles/llm-sizes#model-shrinking))
* quantization results in an quantization error, the less the better

> Quantization: Reducing the precision of weights from floating-point numbers
> (such as, 32-bit) to lower-bit representations (such as, 8-bit).

Common data types:

* FP32 (full precision)
* FP16 (half precision)
* BF16 (a shortened 16-bit version of 32-bit floats)

> It preserves the approximate dynamic range of 32-bit floating-point numbers
> by retaining 8 exponent bits, but supports only an 8-bit precision rather
> than the 24-bit significand of the binary32 format. -- [bfloat16 floating-point format](https://en.wikipedia.org/wiki/Bfloat16_floating-point_format)

> Our results show that deep learning training using BFLOAT16 tensors achieves
> the same state-of-the-art (SOTA) results across domains as FP32 tensors in
> the same number of iterations and with no changes to hyper-parameters. -- [A Study of BFLOAT16 for Deep Learning Training](https://arxiv.org/pdf/1905.12322)

* INT8 (1 byte)

### Suffixes

* https://huggingface.co/docs/hub/en/gguf#quantization-types
* TODO: https://medium.com/@paul.ilvez/demystifying-llm-quantization-suffixes-what-q4-k-m-q8-0-and-q6-k-really-mean-0ec2770f17d3

Deep dive: [A Visual Guide to Quantization](https://www.maartengrootendorst.com/blog/quantization/)


