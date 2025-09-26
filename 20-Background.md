# Background

* interest in NLP since early 2000s,
  [wortschatz](https://wortschatz-leipzig.de/en)
([ex](https://dict.wortschatz-leipzig.de/en/res?corpusId=eng_news_2024&word=Uffizi+Gallery))
at Leipzig University

![Picture of a book: Dornseiff dictionary](static/9783110002874-de.jpg)

> 1933/34 unter dem Titel *Der deutsche Wortschatz synonymisch geordnet.* --
> [Dornseiff - Der deutsche Wortschatz nach
> Sachgruppen](https://ids-pub.bsz-bw.de/frontdoor/deliver/index/docId/4961/file/Storjohann_Dornseiff_Der_deutsche_Wortschatz_nach_Sachgruppen_2012.pdf)

![](static/dornseiff-page.png)

Similar, popular, early project: [wordnet](https://en.wikipedia.org/wiki/WordNet):

> WordNet is a lexical database of semantic relations between words that links
> words into semantic relations including synonyms, hyponyms, and meronyms.

```sh
$ sudo apt install wordnet
$ wnb
```

![WordNet Browser Overview: go](static/screenshot-2025-09-20-170918-wordnet-go.png)

![WordNet Browser Hyponyms: train](static/screenshot-2025-09-20-170947-wordnet-train-hyponyms.png)

Statistical approach, counting words, sparse representation (bag-of-words), manual curation.

![](static/screenshot-2025-09-20-171719-gemini-norvig-quote.png)

Back then, there was an influencial paper about how more data wins over
algorithms, let me quickly ask an LLM which paper it was?

![](static/screenshot-2025-09-20-172118-claude-scaling-2001.png)

A few open weights LLMs struggle a bit with this question, but [Meta Llama 3.3
70B Instruct](https://huggingface.co/meta-llama/Llama-3.3-70B-Instruct)
(released December 6, 2024) seems to remember it.

![](static/screenshot-2025-09-20-172808-chatai-meta-llama3.3-70B-scaling-2001.png)

The paper was called: "Scaling to Very Very Large Corpora for Natural Language
Disambiguation" (2001)

> We collected a **1-billion-word** training corpus from a variety of English
> texts, including news articles, scientific abstracts, government transcripts,
> literature and other varied forms of prose. This training corpus is three

