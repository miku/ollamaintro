"""
Chunk Sherlock stories into smaller chunks.
"""

import fileinput
import json
import os

def main():
    snippets = []
    batch = []
    for line in fileinput.input():
        line = line.strip()
        if len(line) == 0:
            if batch:
                text = " ".join(batch)
                snippets.append(text)
                batch = []
        else:
            batch.append(line)
    if batch:
        snippets.append(batch)


    if os.environ.get("DEBUG"):
        print(len(snippets))
        for i, s in enumerate(snippets):
            print()
            print(f" ==== {i} ====")
            print(f"{s}")
    else:
        for i, s in enumerate(snippets):
            print(json.dumps({"id": i, "text": s, "len": len(s)}))

if __name__ == '__main__':
    main()
