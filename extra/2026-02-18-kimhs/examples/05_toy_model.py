#!/usr/bin/env python3
"""
05_toy_model.py

Create a minimal toy model with random weights in GGUF format.
This demonstrates how to build a custom architecture that Ollama can load.

Requirements:
    pip install gguf numpy

Run:
    python 05_toy_model.py --hidden-dim 512 --n-heads 8 --vocab-size 32000

Testing with ollama runner directly:
    # Start the runner
    ollama runner --ollama-engine --model ./toymodel.gguf --port 8080 &

    # Load the model (must specify parameters)
    curl -X POST localhost:8080/load -d '{
        "Operation": 2,
        "Parallel": 1,
        "BatchSize": 512,
        "KvSize": 2048
    }'

    # Send a completion request
    curl -X POST localhost:8080/completion -d '{"prompt": "Hello", "n_predict": 10}'

Note: The output will be random since weights are untrained!
The model demonstrates valid GGUF structure compatible with ollama's engine.
"""

import numpy as np
import struct

# Import gguf - pip install gguf
try:
    import gguf
    from gguf import GGUFWriter, GGMLQuantizationType
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
    - blk.N.ffn_up/down/gate.weight: feed-forward network
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
    writer = GGUFWriter(output_path, "llama")

    # Set random seed for reproducibility
    np.random.seed(42)

    head_dim = hidden_dim // n_heads
    ffn_dim = hidden_dim * 4  # Standard 4x expansion

    # ====================
    # Add required metadata
    # ====================

    # General metadata
    writer.add_name("ToyModel")
    writer.add_type("model")  # Required: indicates this is a model
    writer.add_description("A minimal toy model for demonstration")
    writer.add_file_type(GGMLQuantizationType.F32)

    # Architecture-specific metadata (llama.*)
    writer.add_context_length(context_length)
    writer.add_embedding_length(hidden_dim)
    writer.add_block_count(n_layers)
    writer.add_head_count(n_heads)
    writer.add_head_count_kv(n_heads)  # Same as head_count for standard MHA
    writer.add_feed_forward_length(ffn_dim)
    writer.add_rope_dimension_count(head_dim)
    writer.add_rope_freq_base(10000.0)
    writer.add_layer_norm_rms_eps(1e-5)
    writer.add_key_length(head_dim)
    writer.add_value_length(head_dim)

    # Vocabulary size
    writer.add_vocab_size(vocab_size)

    print("Core metadata added.")

    # ====================
    # Add tokenizer metadata
    # ====================

    # Use GPT-2 style BPE tokenizer
    writer.add_tokenizer_model("gpt2")
    writer.add_tokenizer_pre("default")  # Preprocessing type

    # Create token list - byte-level tokens
    tokens = []
    for i in range(vocab_size):
        if i == 0:
            tokens.append("<pad>")
        elif i == 1:
            tokens.append("<s>")  # BOS
        elif i == 2:
            tokens.append("</s>")  # EOS
        elif i == 3:
            tokens.append("<unk>")
        elif i < 128:
            # Printable ASCII
            if 32 <= i < 127:
                tokens.append(chr(i))
            else:
                tokens.append(f"<0x{i:02X}>")
        else:
            tokens.append(f"<0x{i:02X}>")

    writer.add_token_list(tokens)

    # Token types: 1=normal, 2=unknown, 3=control, 4=user_defined, 5=unused, 6=byte
    token_types = []
    for i in range(vocab_size):
        if i == 0:
            token_types.append(3)  # control (pad)
        elif i == 1:
            token_types.append(3)  # control (bos)
        elif i == 2:
            token_types.append(3)  # control (eos)
        elif i == 3:
            token_types.append(2)  # unknown
        else:
            token_types.append(1)  # normal

    writer.add_token_types(token_types)

    # Token scores (for SentencePiece compatibility, but needed)
    scores = [0.0] * vocab_size
    writer.add_token_scores(scores)

    # Special tokens
    writer.add_bos_token_id(1)
    writer.add_eos_token_id(2)
    writer.add_pad_token_id(0)
    writer.add_add_bos_token(True)
    writer.add_add_eos_token(False)

    # BPE merges - minimal set (empty is OK for byte-level)
    # For a real model, you'd have thousands of merges
    merges = []
    writer.add_token_merges(merges)

    print("Tokenizer metadata added.")

    # ====================
    # Add model tensors
    # ====================

    # Token embeddings: [vocab_size, hidden_dim]
    add_tensor(writer, "token_embd.weight",
               shape=(vocab_size, hidden_dim))

    # Transformer blocks
    for layer in range(n_layers):
        prefix = f"blk.{layer}"

        # Attention weights
        # Q, K, V projections
        add_tensor(writer, f"{prefix}.attn_q.weight",
                   shape=(hidden_dim, hidden_dim))
        add_tensor(writer, f"{prefix}.attn_k.weight",
                   shape=(hidden_dim, hidden_dim))
        add_tensor(writer, f"{prefix}.attn_v.weight",
                   shape=(hidden_dim, hidden_dim))
        add_tensor(writer, f"{prefix}.attn_output.weight",
                   shape=(hidden_dim, hidden_dim))

        # Feed-forward weights (SwiGLU style with gate)
        add_tensor(writer, f"{prefix}.ffn_gate.weight",
                   shape=(ffn_dim, hidden_dim))
        add_tensor(writer, f"{prefix}.ffn_up.weight",
                   shape=(ffn_dim, hidden_dim))
        add_tensor(writer, f"{prefix}.ffn_down.weight",
                   shape=(hidden_dim, ffn_dim))

        # RMSNorm weights (1D, no bias for RMSNorm)
        add_tensor(writer, f"{prefix}.attn_norm.weight",
                   shape=(hidden_dim,), init_ones=True)
        add_tensor(writer, f"{prefix}.ffn_norm.weight",
                   shape=(hidden_dim,), init_ones=True)

    # Output layers
    add_tensor(writer, "output_norm.weight",
               shape=(hidden_dim,), init_ones=True)
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
    print()
    print("  # Create a Modelfile that references the GGUF")
    print(f"  echo 'FROM ./{output_path}' > Modelfile")
    print("  ollama create toymodel -f Modelfile")
    print("  ollama run toymodel 'Hello'")
    print()
    print("Note: Output will be random (untrained weights)!")

    # Also write the Modelfile for convenience
    modelfile_path = output_path.replace('.gguf', '.Modelfile')
    with open(modelfile_path, 'w') as f:
        f.write(f"FROM ./{output_path}\n")
        f.write("PARAMETER temperature 0.8\n")
        f.write("PARAMETER num_ctx 512\n")
    print(f"Modelfile written to: {modelfile_path}")


def add_tensor(writer, name: str, shape: tuple, init_ones: bool = False):
    """Add a tensor with random or ones initialization."""
    if init_ones:
        # For normalization layers, initialize to ones
        data = np.ones(shape, dtype=np.float32)
    else:
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
