import pandas as pd
import requests
import json
import time
import datetime


# https://medium.com/@RareLoot/using-pushshifts-api-to-extract-reddit-submissions-fb517b286563

def getPushshiftData(query, after, before, sub):
    url = 'https://api.pushshift.io/reddit/search/submission/?title=' + str(query) + '&size=1000&after=' + str(
        after) + '&before=' + str(before) + '&subreddit=' + str(sub)
    print(url)
    r = requests.get(url)
    data = json.loads(r.text)
    cleanedData = []
    for i in range(len(data['data'])):
        if data['data'][i]['author'] != '[deleted]' and data['data'][i]['score'] > 0:
            print(data['data'][i])
            cleanedData.append(data['data'][i])
    return cleanedData


def collectSubData(subm):
    subData = list()  # list to store data points
    title = subm['title']
    url = subm['url']
    try:
        flair = subm['link_flair_text']
    except KeyError:
        flair = "NaN"
    author = subm['author']
    sub_id = subm['id']
    score = subm['score']
    created = datetime.datetime.fromtimestamp(subm['created_utc'])  # 1520561700.0
    numComms = subm['num_comments']
    permalink = subm['permalink']

    subData.append((sub_id, title, url, author, score, created, numComms, permalink, flair))
    subStats[sub_id] = subData


# Subreddit to query
sub = 'funny'
# before and after dates
before = "1317470400"
after = "1285934400"
query = ""
subCount = 0
subStats = {}

data = getPushshiftData(query, after, before, sub)  # Will run until all posts have been gathered
# from the 'after' date up until before date
while len(data) > 0:
    for submission in data:
        collectSubData(submission)
        subCount += 1
    # Calls getPushshiftData() with the created date of the last submission
    print(len(data))
    print(str(datetime.datetime.fromtimestamp(data[-1]['created_utc'])))
    after = data[-1]['created_utc']
    data = getPushshiftData(query, after, before, sub)
    #print(data)
    # add data to DB?

print(len(data))
