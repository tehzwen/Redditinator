import sys
from textblob import TextBlob

def main():

    if (len(sys.argv) < 1):
        exit(-1)

    text = TextBlob(sys.argv[1])
    val = {
        "polarity" : text.polarity,
        "subjectivity" : text.subjectivity,
        "tags": text.tags,
        
    }

    print(val)

main()
