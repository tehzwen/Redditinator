import sys
from textblob import TextBlob

nltk.download('punkt')

def main():

    #pull out args
    args = str(sys.argv)

    if (len(args) < 0):
        exit(-1)

    text = TextBlob("Hey how are you?")
    print("TEXTBLOB: ", text.sentiment)

main()
