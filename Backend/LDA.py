import nltk
from nltk import FreqDist
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

# Initialize spacy 'en' model, keeping only tagger component (for efficiency)
nlp = spacy.load("en_core_web_sm", disable=['parser', 'ner'])

# Referencing tutorial: https://www.machinelearningplus.com/nlp/topic-modeling-gensim-python/
RUNUPDATE = True
def sent_to_words(sentences):
    for sentence in sentences:
        yield(gensim.utils.simple_preprocess(str(sentence), deacc=True))  # deacc=True removes punctuations

def remove_stopwords(texts):
    return [[word for word in simple_preprocess(str(doc)) if word not in stop_words] for doc in texts]



def lemmatization(texts, allowed_postags=['NOUN', 'ADJ', 'VERB', 'ADV']):
    """https://spacy.io/api/annotation"""
    texts_out = []
    for sent in texts:
        doc = nlp(" ".join(sent))
        texts_out.append([token.lemma_ for token in doc if token.pos_ in allowed_postags])
    return texts_out


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
   
        headers = {
        'Content-Type': 'text/plain'
        }
       
        response = requests.request("GET", url, headers=headers)

        data = response.json()
        return data

def updateTopic(postID, topic):
        url = "http://localhost:4000/topic"

        payload = "{ \n\"id\":\""+postID+ "\",\n\"topic\":\""+topic+"\" \n}"
        headers = {
        'Content-Type': 'text/plain'
        }
       
        response = requests.request("POST", url, headers=headers, data = payload)
        return response

def run(subreddit):
        posts = getPosts(subreddit)        

        data = []
        for post in posts:
                if post['num_comments'] != '0':
                        postId = post['id']
                        
                postComments = []
                postComments.append(post['title'])
                        
                if post['selftext'] != '':
                        postComments.append(post['selftext'])  

                res = getComments(postId)
                for val in res:
                        temp = val['body']
                        postComments.append(temp)
                        data.append(temp)
                
                allWords = [word for (word, pos) in nltk.pos_tag(nltk.word_tokenize(" ".join(postComments))) if pos[0] == 'N']
                #allWords = removeSymbols(allWords)
                fdist = FreqDist(allWords)
                commonNoun = fdist.most_common(1)
                if (len(commonNoun) != 0 and commonNoun[0][0].isalpha() and RUNUPDATE):
                        updateTopic(postId, commonNoun[0][0].lower())
                
        print("Updated post topics! Analyzing subreddit topics.....")

        # Remove new line characters
        data = [re.sub('\s+', ' ', sent) for sent in data]

        # Remove distracting single quotes
        data = [re.sub("\'", "", sent) for sent in data]

        # Define functions for stopwords, bigrams, trigrams and lemmatization

        # Creates list of cleaned documents
        data_words = list(sent_to_words(data))

        # Build the bigram and trigram models
        bigram = gensim.models.Phrases(data_words, min_count=5, threshold=100) # higher threshold fewer phrases.
        trigram = gensim.models.Phrases(bigram[data_words], threshold=100)
        def make_bigrams(texts):
                return [bigram_mod[doc] for doc in texts]

        def make_trigrams(texts):
                return [trigram_mod[bigram_mod[doc]] for doc in texts]
        # Faster way to get a sentence clubbed as a trigram/bigram
        bigram_mod = gensim.models.phrases.Phraser(bigram)
        trigram_mod = gensim.models.phrases.Phraser(trigram)

        # Remove Stop Words
        data_words_nostops = remove_stopwords(data_words)

        # Form Bigrams
        data_words_bigrams = make_bigrams(data_words_nostops)



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
        # vis = pyLDAvis.gensim.prepare(lda_model, corpus, id2word)
        # pyLDAvis.show(vis)
subreddits = [
"calgary",
"democrats",
"republican",
"alberta",
"exmormon",
"Edmonton",
"witcher",
"financialindependence",
"tifu",
"democrats",
"wholesomememes"]

for sub in subreddits:
        run(sub)
