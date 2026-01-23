#!/usr/bin/env python3
"""
05_toy_model.py

Create a minimal toy model with random weights in GGUF format.
This demonstrates how to build a custom architecture that Ollama can load.

Requirements:
    pip install gguf numpy

Run:
    python 05_toy_model.py

This creates 'toymodel.gguf' which can be imported into Ollama:
    ollama create toymodel -f toymodel.gguf
    ollama run toymodel "Hello"

Note: The output will be random since weights are random!
"""

import numpy as np

# Import gguf - pip install gguf
try:
    import gguf
except ImportError:
    print("Please install gguf: pip install gguf")
    exit(1)


def create_toy_model(
    output_path: str = "toymodel.gguf",
    vocab_size: int = 256,      # Small vocabulary (byte-level)
    hidden_dim: int = 64,       # Hidden dimension
    n_layers: int = 2,          # Number of transformer blocks
    n_heads: int = 4,           # Number of attention heads
    context_length: int = 512,  # Maximum context length
):
    """
    Create a minimal GGUF model file with random weights.

    The architecture follows the standard llama pattern:
    - token_embd.weight: vocabulary embeddings
    - blk.N.attn_q/k/v.weight: attention projections
    - blk.N.attn_output.weight: attention output projection
    - blk.N.ffn_up/down.weight: feed-forward network
    - blk.N.attn_norm/ffn_norm.weight: layer norms
    - output_norm.weight: final layer norm
    - output.weight: output projection to vocabulary
    """

    print(f"Creating toy model: {output_path}")
    print(f"  Vocabulary size: {vocab_size}")
    print(f"  Hidden dimension: {hidden_dim}")
    print(f"  Number of layers: {n_layers}")
    print(f"  Attention heads: {n_heads}")
    print(f"  Context length: {context_length}")
    print()

    # Initialize GGUF writer with "llama" architecture
    # This tells Ollama to use the llama model implementation
    writer = gguf.GGUFWriter(output_path, "llama")

    # Set random seed for reproducibility
    np.random.seed(42)

    # ====================
    # Add required metadata
    # ====================

    # General metadata
    writer.add_name("ToyModel")
    writer.add_description("A minimal toy model for demonstration")
    writer.add_file_type(gguf.GGMLQuantizationType.F32)  # Full precision

    # Architecture-specific metadata (llama.*)
    writer.add_context_length(context_length)
    writer.add_embedding_length(hidden_dim)
    writer.add_block_count(n_layers)
    writer.add_head_count(n_heads)
    writer.add_head_count_kv(n_heads)  # For grouped-query attention (same as head_count for MHA)
    writer.add_feed_forward_length(hidden_dim * 4)  # Standard 4x expansion
    writer.add_rope_dimension_count(hidden_dim // n_heads)  # RoPE dimensions

    # Tokenizer metadata (minimal byte-level tokenizer)
    writer.add_tokenizer_model("gpt2")  # Use GPT-2 style tokenizer
    writer.add_vocab_size(vocab_size)

    # Add token list (byte values as strings for simplicity)
    tokens = [f"<{i:02x}>" for i in range(vocab_size)]
    writer.add_token_list(tokens)

    # Special tokens
    writer.add_bos_token_id(1)
    writer.add_eos_token_id(2)
    writer.add_pad_token_id(0)

    print("Metadata added.")

    # ====================
    # Add model tensors
    # ====================

    head_dim = hidden_dim // n_heads
    ffn_dim = hidden_dim * 4

    # Token embeddings: [vocab_size, hidden_dim]
    add_tensor(writer, "token_embd.weight",
               shape=(vocab_size, hidden_dim))

    # Transformer blocks
    for layer in range(n_layers):
        prefix = f"blk.{layer}"

        # Attention weights
        add_tensor(writer, f"{prefix}.attn_q.weight",
                   shape=(hidden_dim, hidden_dim))
        add_tensor(writer, f"{prefix}.attn_k.weight",
                   shape=(hidden_dim, hidden_dim))
        add_tensor(writer, f"{prefix}.attn_v.weight",
                   shape=(hidden_dim, hidden_dim))
        add_tensor(writer, f"{prefix}.attn_output.weight",
                   shape=(hidden_dim, hidden_dim))

        # Feed-forward weights
        add_tensor(writer, f"{prefix}.ffn_up.weight",
                   shape=(ffn_dim, hidden_dim))
        add_tensor(writer, f"{prefix}.ffn_down.weight",
                   shape=(hidden_dim, ffn_dim))

        # Layer norms (1D)
        add_tensor(writer, f"{prefix}.attn_norm.weight",
                   shape=(hidden_dim,))
        add_tensor(writer, f"{prefix}.ffn_norm.weight",
                   shape=(hidden_dim,))

    # Output layers
    add_tensor(writer, "output_norm.weight",
               shape=(hidden_dim,))
    add_tensor(writer, "output.weight",
               shape=(vocab_size, hidden_dim))

    print("Tensors added.")

    # ====================
    # Write the file
    # ====================

    writer.write_header_to_file()
    writer.write_kv_data_to_file()
    writer.write_tensors_to_file()
    writer.close()

    print(f"\nModel written to: {output_path}")
    print("\nTo use with Ollama:")
    print(f"  ollama create toymodel -f {output_path}")
    print("  ollama run toymodel 'Hello'")
    print("\nNote: Output will be random (untrained weights)!")


def add_tensor(writer, name: str, shape: tuple):
    """Add a tensor with random weights (Xavier initialization)."""
    # Xavier/Glorot initialization for better random outputs
    fan_in = shape[-1] if len(shape) > 1 else shape[0]
    fan_out = shape[0] if len(shape) > 1 else shape[0]
    std = np.sqrt(2.0 / (fan_in + fan_out))

    data = np.random.normal(0, std, shape).astype(np.float32)

    print(f"  Adding: {name:40s} shape={shape}")
    writer.add_tensor(name, data)


if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser(description="Create a toy GGUF model")
    parser.add_argument("-o", "--output", default="toymodel.gguf",
                        help="Output file path")
    parser.add_argument("--vocab-size", type=int, default=256,
                        help="Vocabulary size")
    parser.add_argument("--hidden-dim", type=int, default=64,
                        help="Hidden dimension")
    parser.add_argument("--n-layers", type=int, default=2,
                        help="Number of transformer layers")
    parser.add_argument("--n-heads", type=int, default=4,
                        help="Number of attention heads")
    parser.add_argument("--context-length", type=int, default=512,
                        help="Maximum context length")

    args = parser.parse_args()

    create_toy_model(
        output_path=args.output,
        vocab_size=args.vocab_size,
        hidden_dim=args.hidden_dim,
        n_layers=args.n_layers,
        n_heads=args.n_heads,
        context_length=args.context_length,
    )
