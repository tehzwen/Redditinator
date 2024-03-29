import nltk
import sys
import time
from nltk.corpus import stopwords
from nltk.tokenize import RegexpTokenizer
from nltk.stem.porter import PorterStemmer
from gensim import corpora, models
import gensim
import requests


def analyzeTopic(value, stopWords):

    tokenizer = RegexpTokenizer(r'\w+')

    # create English stop words list

    # Create p_stemmer of class PorterStemmer
    p_stemmer = PorterStemmer()

    # list for tokenized documents in loop
    texts = []

    # Convert strings to raw text
    raw = value.lower()
    tokens = tokenizer.tokenize(raw)

    # remove stop words from tokens
    stopped_tokens = [i for i in tokens if not i in stopWords]

    # stem tokens
    stemmed_tokens = [p_stemmer.stem(i) for i in stopped_tokens]

    # add tokens to list
    texts.append(stemmed_tokens)
    # turn our tokenized documents into a id <-> term dictionary
    dictionary = corpora.Dictionary(texts)

    # convert tokenized documents into a document-term matrix
    corpus = [dictionary.doc2bow(text) for text in texts]

    # generate LDA model
    ldamodel = gensim.models.ldamodel.LdaModel(
        corpus, num_topics=1, id2word=dictionary, passes=20)

    print(ldamodel.print_topics(num_topics=1, num_words=1))


def main():
    nltk.download('stopwords', quiet=True)
    nStopWords = stopwords.words('english')

    while(True):
        time.sleep(0.2)
        data = sys.stdin.readlines()
        if (len(data) > 0):
            analyzeTopic(data[0], nStopWords)

main()
