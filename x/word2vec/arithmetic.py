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
    # king - man + women = ...

    print("\033[0;32mNIGHT-DARK\033[0m")
    for w in model.wv.most_similar(positive=["night"], negative=["dark"], topn=3):
        print("\t".join((str(v) for v in w)))
    # print()
    # print("\033[0;32mWAR-BLOOD\033[0m")
    # for w in model.wv.most_similar(positive=["", ""], negative=[], topn=3):
    #     print("\t".join((str(v) for v in w)))
    print()
    print("\033[0;32mSHIP+MOVE\033[0m")
    for w in model.wv.most_similar(positive=["ship", "move"], topn=3):
        print("\t".join((str(v) for v in w)))
    print()
    print("\033[0;32mCITY+PEOPLE\033[0m")
    for w in model.wv.most_similar(positive=["city", "people"], topn=3):
        print("\t".join((str(v) for v in w)))

    # $ python arithmetic.py | column -t
    # KIND+WOMEN-MAN
    # queen           0.6779043674468994
    # captain         0.6737080216407776
    # prophet         0.6566919684410095
    # WAR-BLOOD
    # 139:20          0.5924164652824402
    # kadeshbarnea    0.566044270992279
    # scuffling       0.5186205506324768
    # SHIP+MOVE
    # sail            0.9151894450187683
    # close           0.878623902797699
    # characters      0.8774296641349792
    # CITY+PEOPLE
    # congregation    0.7869210243225098
    # camp            0.778090238571167
    # land            0.7592784762382507

