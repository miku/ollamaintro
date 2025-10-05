import sys
from gensim.models import Word2Vec

MODELFILE = "word2vec.model"

if __name__ == '__main__':
    if len(sys.argv) > 1:
        word = sys.argv[1].lower()
    else:
        word = "computer"
    model = Word2Vec.load(MODELFILE)
    for word, dist in model.wv.most_similar(word):
        print(word, dist)

# $ python most_similar.py woman | column -t
# child      0.8833639025688171
# lady       0.7483692169189453
# deceiver   0.7366161346435547
# fool       0.7194058895111084
# toy        0.7111260294914246
# sign       0.7083047032356262
# girl       0.7038264870643616
# friend     0.7031304836273193
# galilaean  0.7015317678451538
# refiner    0.69990074634552

# $ python most_similar.py man | column -t
# thing          0.7021642923355103
# woman          0.678594708442688
# blighted       0.6671753525733948
# nerve          0.6537002325057983
# one            0.6505756378173828
# fool           0.6469791531562805
# penny          0.6439594030380249
# well-educated  0.6354129314422607
# child          0.6281625032424927
# person         0.6188145279884338

# $ python most_similar.py library | column -t
# copy         0.839066743850708
# memory       0.8309504985809326
# poem         0.8284162878990173
# reasonable   0.8268866539001465
# alarming     0.8206095695495605
# description  0.8205781579017639
# taste        0.8205443024635315
# lack         0.8198086023330688
# meaning      0.8157389760017395
# exporting    0.8156920075416565

