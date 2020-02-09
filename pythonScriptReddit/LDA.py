import nltk
from nltk.corpus import stopwords
from nltk.tokenize import RegexpTokenizer
from nltk.stem.porter import PorterStemmer
from gensim import corpora, models
import gensim
import requests

tokenizer = RegexpTokenizer(r'\w+')

# create English stop words list
nltk.download('stopwords')
stopwords = stopwords.words('english')

# Create p_stemmer of class PorterStemmer
p_stemmer = PorterStemmer()

# create sample documents
# long = requests.get('http://rare-technologies.com/the_matrix_synopsis.txt').text
long = (
    "I got my dog after a suicide attempt that landed me in the hospital for a while. I was beyond depressed, "
    "abusing alcohol, and I just had a thought that I wanted a dog.I was scrolling through adoption ads, going into "
    "random pet stores, checking all the rescues. And then one day I saw my pup in a little cage with some other "
    "dogs. He wasn’t even three pounds and he didn’t give a FUCK that anyone was looking at him. He was just chilling "
    "and doing his own thing, all the other puppies were jumping around.I held him in my arms and I just had this "
    "warm feeling fill my entire chest. It was the first time I’d really felt like that in a long time. It’s been 2.5 "
    "years now and I’m not depressed anymore. My life has done a complete 180 and my dog is currently sitting here "
    "grabbing my hand like “excuse me you’re welcome now pet me”. Any time the thought that life isn’t going my way "
    "comes, I look at him and remember that everything has to be okay bc he needs me.So I agree, the world’s best "
    "antidepressant! "
)

short = "Donald Trump does not have the values of America. he continuously damages our way of life. "
# list for tokenized documents in loop
texts = []

# Convert strings to raw text
raw = short.lower()
tokens = tokenizer.tokenize(raw)

# remove stop words from tokens
stopped_tokens = [i for i in tokens if not i in stopwords]

# stem tokens
stemmed_tokens = [p_stemmer.stem(i) for i in stopped_tokens]

# add tokens to list
texts.append(stemmed_tokens)
# turn our tokenized documents into a id <-> term dictionary
dictionary = corpora.Dictionary(texts)

# convert tokenized documents into a document-term matrix
corpus = [dictionary.doc2bow(text) for text in texts]

# generate LDA model
ldamodel = gensim.models.ldamodel.LdaModel(corpus, num_topics=1, id2word=dictionary, passes=20)

print(ldamodel.print_topics(num_topics=1, num_words=1))
