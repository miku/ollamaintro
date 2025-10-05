import sys
import itertools
from gensim.models import Word2Vec

MODELFILE = "word2vec.model"

if __name__ == '__main__':
    if len(sys.argv) > 1:
        word = sys.argv[1].lower()
    else:
        word = "computer"
    model = Word2Vec.load(MODELFILE)
    words = ["cat", "mouse", "table"]
    for a, b in itertools.combinations(words, r=2):
        # for d in model.wv.distances(a, [a]):
        #     print(a, a, abs(round(d, 4)))
        d = model.wv.distance(a, b)
        print(a, b, abs(round(d, 4)))

        # $ python dist.py | column -t | sort -k 3,3 -n
        # cat    mouse  0.2123
        # cat    table  0.3126
        # mouse  table  0.4248

