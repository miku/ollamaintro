#!/usr/bin/env -S uv run --script
# /// script
# requires-python = ">=3.12"
# dependencies = ["numpy"]
# ///
import numpy as np
from collections import Counter, defaultdict

class Word2Vec:
    def __init__(self, embedding_dim=100, window_size=5, min_count=5,
                 negative_samples=5, learning_rate=0.025, epochs=5):
        """
        Simple Word2Vec implementation with Skip-gram and Negative Sampling

        Args:
            embedding_dim: Dimension of word embeddings
            window_size: Context window size
            min_count: Minimum word frequency to include in vocabulary
            negative_samples: Number of negative samples per positive sample
            learning_rate: Learning rate for SGD
            epochs: Number of training epochs
        """
        self.embedding_dim = embedding_dim
        self.window_size = window_size
        self.min_count = min_count
        self.negative_samples = negative_samples
        self.learning_rate = learning_rate
        self.epochs = epochs

        self.vocab = {}  # word -> index
        self.index_to_word = {}  # index -> word
        self.word_freq = None
        self.W_input = None  # Input word embeddings
        self.W_output = None  # Output word embeddings

    def _build_vocab(self, sentences):
        """Build vocabulary from sentences"""
        word_counts = Counter()
        for sentence in sentences:
            word_counts.update(sentence)

        # Filter by min_count
        vocab_words = [word for word, count in word_counts.items()
                       if count >= self.min_count]

        # Create mappings
        self.vocab = {word: idx for idx, word in enumerate(vocab_words)}
        self.index_to_word = {idx: word for word, idx in self.vocab.items()}

        # Store word frequencies for negative sampling
        self.word_freq = np.array([word_counts[word] for word in vocab_words])
        # Raise to power 0.75 (empirically better for negative sampling)
        self.word_freq = np.power(self.word_freq, 0.75)
        self.word_freq = self.word_freq / self.word_freq.sum()

        return len(self.vocab)

    def _get_negative_samples(self, target_idx, n_samples):
        """Sample negative examples (words not in context)"""
        negative_samples = []
        while len(negative_samples) < n_samples:
            sample = np.random.choice(len(self.vocab), p=self.word_freq)
            if sample != target_idx:
                negative_samples.append(sample)
        return negative_samples

    def _sigmoid(self, x):
        """Sigmoid activation function"""
        return 1 / (1 + np.exp(-np.clip(x, -500, 500)))

    def _train_pair(self, center_idx, context_idx):
        """Train on a single (center, context) pair with negative sampling"""
        # Positive sample
        h = self.W_input[center_idx]  # Hidden layer (center word embedding)
        u = self.W_output[context_idx]  # Output layer for context word

        # Compute gradient for positive sample
        # ====================================
        # h is the embedding of the center word (from input layer)
        # u is the embedding of the context word (from output layer)
        # The dot product measures how "aligned" these vectors are
        # Higher score = model thinks these words are related
        score = np.dot(h, u)

        # Sigmoid squashes the score to a probability between 0 and 1
        # What's the probability this context word appears near the center word?
        pred = self._sigmoid(score)


        grad = (pred - 1) * self.learning_rate  # Error signal

        # Update embeddings
        self.W_output[context_idx] -= grad * h
        grad_h = grad * u

        # Negative samples
        negative_indices = self._get_negative_samples(context_idx, self.negative_samples)
        for neg_idx in negative_indices:
            u_neg = self.W_output[neg_idx]
            score_neg = np.dot(h, u_neg)
            pred_neg = self._sigmoid(score_neg)
            grad_neg = pred_neg * self.learning_rate

            self.W_output[neg_idx] -= grad_neg * h
            grad_h += grad_neg * u_neg

        self.W_input[center_idx] -= grad_h

    def train(self, sentences):
        """
        Train Word2Vec model

        Args:
            sentences: List of sentences, where each sentence is a list of words
        """
        # Build vocabulary
        vocab_size = self._build_vocab(sentences)
        print(f"Vocabulary size: {vocab_size}")

        # Initialize embeddings with small random values
        self.W_input = np.random.uniform(-0.5, 0.5,
                                        (vocab_size, self.embedding_dim)) / self.embedding_dim
        self.W_output = np.random.uniform(-0.5, 0.5,
                                         (vocab_size, self.embedding_dim)) / self.embedding_dim

        # Training loop
        for epoch in range(self.epochs):
            total_pairs = 0

            for sentence in sentences:
                # Convert words to indices
                word_indices = [self.vocab[word] for word in sentence
                              if word in self.vocab]

                # Generate training pairs
                for i, center_idx in enumerate(word_indices):
                    # Get context window
                    start = max(0, i - self.window_size)
                    end = min(len(word_indices), i + self.window_size + 1)

                    for j in range(start, end):
                        if i != j:
                            context_idx = word_indices[j]
                            self._train_pair(center_idx, context_idx)
                            total_pairs += 1

            print(f"Epoch {epoch + 1}/{self.epochs} - Trained on {total_pairs} pairs")

    def get_vector(self, word):
        """Get embedding vector for a word"""
        if word not in self.vocab:
            raise ValueError(f"Word '{word}' not in vocabulary")
        return self.W_input[self.vocab[word]]

    def most_similar(self, word, top_n=10):
        """Find most similar words using cosine similarity"""
        if word not in self.vocab:
            raise ValueError(f"Word '{word}' not in vocabulary")

        word_vec = self.get_vector(word)

        # Compute cosine similarity with all words
        similarities = np.dot(self.W_input, word_vec) / (
            np.linalg.norm(self.W_input, axis=1) * np.linalg.norm(word_vec)
        )

        # Get top N similar words (excluding the word itself)
        similar_indices = np.argsort(similarities)[::-1][1:top_n+1]

        return [(self.index_to_word[idx], similarities[idx])
                for idx in similar_indices]


# Example usage
if __name__ == "__main__":
    # Sample corpus
    sentences = [
        ["the", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"],
        ["the", "dog", "is", "lazy"],
        ["the", "cat", "is", "quick"],
        ["the", "fox", "is", "brown"],
        ["a", "quick", "brown", "dog", "jumps"],
        ["the", "lazy", "cat", "sleeps"],
        ["quick", "animals", "run", "fast"],
        ["brown", "animals", "are", "common"],
    ]

    # Replicate sentences to have more training data
    sentences = sentences * 100

    # Train model
    model = Word2Vec(embedding_dim=50, window_size=2, min_count=1,
                     negative_samples=5, learning_rate=0.025, epochs=10)
    model.train(sentences)

    # Test similarity
    print("\nMost similar to 'quick':")
    for word, similarity in model.most_similar('quick', top_n=5):
        print(f"  {word}: {similarity:.4f}")

    print("\nMost similar to 'dog':")
    for word, similarity in model.most_similar('dog', top_n=5):
        print(f"  {word}: {similarity:.4f}")

# $ ./simple.py
# Installed 1 package in 49ms
# Vocabulary size: 17
# Epoch 1/10 - Trained on 10400 pairs
# Epoch 2/10 - Trained on 10400 pairs
# Epoch 3/10 - Trained on 10400 pairs
# Epoch 4/10 - Trained on 10400 pairs
# Epoch 5/10 - Trained on 10400 pairs
# Epoch 6/10 - Trained on 10400 pairs
# Epoch 7/10 - Trained on 10400 pairs
# Epoch 8/10 - Trained on 10400 pairs
# Epoch 9/10 - Trained on 10400 pairs
# Epoch 10/10 - Trained on 10400 pairs
#
# Most similar to 'quick':
#   fast: 0.3400
#   lazy: 0.3290
#   jumps: 0.3281
#   sleeps: 0.2766
#   are: 0.1895
#
# Most similar to 'dog':
#   fox: 0.6552
#   cat: 0.5560
#   over: 0.4695
#   a: 0.3907
#   is: 0.3636
#
