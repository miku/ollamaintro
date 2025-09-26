# Quantization Notes

* quanto; pytorch model
* model compression

Options:

* model pruning (removing layers that have little effect on performance)
* knowledge distillation (student model; teacher model)

You can quantize weights, or activations.

Common example: FP32 to INT8; using 25% of the storage; aim: reduce
quantization error

## Example

We take float values and map them to integers. Sampling is a form of
quantization. Smaller models, faster inference.

* DeepSeek-R1 is about 720GB in 671B parameters
