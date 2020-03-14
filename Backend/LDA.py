import nltk
import re
import numpy as np
import pandas as pd
from pprint import pprint

import gensim
import gensim.corpora as corpora
from gensim.utils import simple_preprocess
from gensim.models import CoherenceModel

import spacy
import pyLDAvis
import pyLDAvis.gensim  # don't skip this
import matplotlib.pyplot as plt

import requests
import json

from nltk.corpus import stopwords
stop_words = stopwords.words('english')
stop_words.extend(['from', 'subject', 're', 'edu', 'use', 'lines']) #temp

# Referencing tutorial: https://www.machinelearningplus.com/nlp/topic-modeling-gensim-python/

def getPosts(subreddit):
        url = "http://localhost:4000/posts?subreddit="+subreddit

        headers = {
        'Content-Type': 'text/plain'
        }

        response = requests.request("GET", url, headers=headers)
        data = response.json()
        return data

def getComments(postID):
        url = "http://localhost:4000/comments?postID="+postID+"&topLevel=true"
        #url = "http://localhost:4000/comments?postID="+postID

        #payload = "{\n\"subreddits\":"+str(subreddits).replace("\'", "\"")+",\n\"to\":"+end+",\n\"from\":"+start+"\n}"
        headers = {
        'Content-Type': 'text/plain'
        }
       
        response = requests.request("GET", url, headers=headers)

        data = response.json()
        return data

def updateTopic(postID, topic):
        url = "localhost:4000/topic"

        payload = "{ \n\"id\":"+postID+ "\n\"topic\":"+topic+" \n}"
        headers = {
        'Content-Type': 'text/plain'
        }

        response = requests.request("POST", url, headers=headers, data = payload)
        return response


posts = getPosts("edmonton")        
postIds = []
for post in posts:
        if post['num_comments'] != '0':
                postIds.append(post['id'])
        
print(postIds)

#res = getComments("fboo5x")
data = []
for postId in postIds:
        res = getComments(postId)

        for val in res:
                #pprint(val['body'])
                temp = val['body']
                data.append(temp)


# Remove new line characters
data = [re.sub('\s+', ' ', sent) for sent in data]

# Remove distracting single quotes
data = [re.sub("\'", "", sent) for sent in data]

# Define functions for stopwords, bigrams, trigrams and lemmatization
def sent_to_words(sentences):
    for sentence in sentences:
        yield(gensim.utils.simple_preprocess(str(sentence), deacc=True))  # deacc=True removes punctuations

def remove_stopwords(texts):
    return [[word for word in simple_preprocess(str(doc)) if word not in stop_words] for doc in texts]

def make_bigrams(texts):
    return [bigram_mod[doc] for doc in texts]

def make_trigrams(texts):
    return [trigram_mod[bigram_mod[doc]] for doc in texts]

def lemmatization(texts, allowed_postags=['NOUN', 'ADJ', 'VERB', 'ADV']):
    """https://spacy.io/api/annotation"""
    texts_out = []
    for sent in texts:
        doc = nlp(" ".join(sent))
        texts_out.append([token.lemma_ for token in doc if token.pos_ in allowed_postags])
    return texts_out

# Creates list of cleaned documents
data_words = list(sent_to_words(data))

# Build the bigram and trigram models
bigram = gensim.models.Phrases(data_words, min_count=5, threshold=100) # higher threshold fewer phrases.
trigram = gensim.models.Phrases(bigram[data_words], threshold=100)

# Faster way to get a sentence clubbed as a trigram/bigram
bigram_mod = gensim.models.phrases.Phraser(bigram)
trigram_mod = gensim.models.phrases.Phraser(trigram)

# Remove Stop Words
data_words_nostops = remove_stopwords(data_words)

# Form Bigrams
data_words_bigrams = make_bigrams(data_words_nostops)

# Initialize spacy 'en' model, keeping only tagger component (for efficiency)
nlp = spacy.load("en_core_web_sm", disable=['parser', 'ner'])

# Do lemmatization keeping only nouns
data_lemmatized = lemmatization(data_words_bigrams)

# Create Dictionary
id2word = corpora.Dictionary(data_lemmatized)

# Create Corpus
texts = data_lemmatized

# Term Document Frequency
corpus = [id2word.doc2bow(text) for text in texts]

lda_model = gensim.models.ldamodel.LdaModel(corpus=corpus,
                                           id2word=id2word,
                                           num_topics=10,
                                           random_state=100,
                                           update_every=1,
                                           chunksize=100,
                                           passes=100,
                                           alpha='auto',
                                           per_word_topics=True)

pprint(lda_model.print_topics())

# Visual representation
vis = pyLDAvis.gensim.prepare(lda_model, corpus, id2word)
pyLDAvis.show(vis)