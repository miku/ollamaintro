import os
import sys
from gensim.models import Word2Vec
from nltk.tokenize import sent_tokenize, word_tokenize

# corpus.file: 1096052 words, vocab: 32922

MODELFILE = "word2vec.model"

if __name__ == '__main__':
    with open("corpus.file") as f:
        if not os.path.exists(MODELFILE):
            # read data
            data = f.read().replace("\n", " ")
            # collect sentences from data
            sentences= []
            for s in sent_tokenize(data):
                sentences.append([w.lower() for w in word_tokenize(s)])
            # train a model
            model = Word2Vec(sentences=sentences, min_count=1, vector_size=100, window=5, workers=8)
            model.save(MODELFILE)
        else:
            print("loading existing model...", file=sys.stderr)
            model = Word2Vec.load(MODELFILE)
        # output vocabulary
        # for k, v in model.wv.key_to_index.items():
        #     print(f"{k}\t{v}")
        #     # works   306
        #     # seen    307
        #     # has     308
        #     # never   309
        #     # tell    310
        #     # four    311
        #     # drink   312
        #     # five    313
        #     # altar   314
        #     # midst   315
        #     # jacob   316
        print(model.wv["cat"])
        print(model.wv["mouse"])
        print(model.wv["table"])

