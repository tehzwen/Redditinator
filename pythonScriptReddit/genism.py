from pprint import pprint as print
from gensim.summarization import summarize, mz_keywords
from gensim.summarization import keywords
import requests

text = (
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

print(keywords(text))
text = requests.get('http://rare-technologies.com/the_matrix_synopsis.txt').text
print(keywords(text))
# text = requests.get('http://rare-technologies.com/the_matrix_synopsis.txt').text
# print(keywords(text, ratio=0.01))
