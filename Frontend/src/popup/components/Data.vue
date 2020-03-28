<template>
  <div class="background">
    <div>
      <img src="/assets/redditinator.png" />
    </div>
    <h1 class="subreddit-title">r/{{this.state.subreddit.name}}</h1>
    <div class="row">
      <div class="col">
        <div class="row">
          <div class="col">
            <button v-on:click="changeSentimentType('post')" class="my-button">Posts</button>
          </div>
          <div class="col">
            <button v-on:click="changeSentimentType('comment')" class="my-button">Comments</button>
          </div>
        </div>
        <div class="row">
          <div class="col" style="margin-right:15px">
            <img
              v-bind:class="sentImageClasses"
              src="/assets/redditGuy.png"
              width="100px"
              height="100px"
            />
            <p
              class="subreddit-title"
            >Overall: {{sentimentType === 'post' ? sentiment.postSent.toFixed(3) : sentiment.commentSent.toFixed(3)}}</p>
          </div>
          <div class="col">
            <div v-if="sentimentType === 'post'">
              <div class="row">
                <p
                  class="sent-text positive-text"
                >{{(sentiment.postSentPos * 100).toFixed(2)}}% Good</p>
              </div>
              <div class="row">
                <p
                  class="sent-text neutral-text"
                >{{(sentiment.postSentNeu * 100).toFixed(2)}}% Neutral</p>
              </div>
              <div class="row">
                <p class="sent-text negative-text">{{(sentiment.postSentNeg * 100).toFixed(2)}}% Bad</p>
              </div>
            </div>
            <div v-if="sentimentType === 'comment'">
              <div class="row">
                <p
                  class="sent-text positive-text"
                >{{(sentiment.commentSentPos * 100).toFixed(2)}}% Good</p>
              </div>
              <div class="row">
                <p
                  class="sent-text neutral-text"
                >{{(sentiment.commentSentNeu * 100).toFixed(2)}}% Neutral</p>
              </div>
              <div class="row">
                <p
                  class="sent-text negative-text"
                >{{(sentiment.commentSentNeg * 100).toFixed(2)}}% Bad</p>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="col" />
    </div>
    <button class="my-button" v-on:click="goToHome">Go to home!</button>
    <MyWordCloud />
  </div>
</template>

<script>
import Axios from "axios";
import MyWordCloud from "./MyWordCloud";

export default {
  props: {
    state: Object
  },
  components: {
    MyWordCloud
  },
  data: () => {
    return {
      topicCounts: [], //contains data for wordcloud
      sentiment: {},
      sentimentType: "post",
      sentImageClasses: ["sentiment-image"]
    };
  },
  mounted() {
    let port = chrome.extension.connect({
      name: "Sample Communication"
    });
    this.port = port;

    //fetch the word cloud data for this subreddit
    Axios.get(
      "http://167.172.132.5:4000/subreddits/" +
        this.state.subreddit.id +
        "/topics"
    )
      .then(res => {
        //this.port.postMessage(res);
        this.topicCounts = res.data;
      })
      .catch(err => {
        this.port.postMessage("Error: ", err);
      });

    Axios.get(
      "http://167.172.132.5:4000/subreddits/" +
        this.state.subreddit.id +
        "/sentiment"
    )
      .then(res => {
        //this.port.postMessage(res);
        this.sentiment = res.data;
      })
      .then(() => {
        this.getImageClass();
      })
      .catch(err => {
        this.port.postMessage("Error: Couldn't fetch sentiments");
      });
  },
  methods: {
    goToHome() {
      this.$emit("page", { page: "home", state: {} });
    },
    changeSentimentType(val) {
      this.sentimentType = val;
      this.getImageClass();
    },
    getImageClass() {
      if (this.sentimentType === "post") {
        if (this.sentiment.postSent > 0) {
          this.sentImageClasses = ["sentiment-image", "positive-background"];
        } else if (this.sentiment.postSent < 0) {
          this.sentImageClasses = ["sentiment-image", "negative-background"];
        } else {
          this.sentImageClasses = ["sentiment-image", "neutral-background"];
        }
      } else {
        if (this.sentiment.commentSent > 0) {
          this.sentImageClasses = ["sentiment-image", "positive-background"];
        } else if (this.sentiment.commentSent < 0) {
          this.sentImageClasses = ["sentiment-image", "negative-background"];
        } else {
          this.sentImageClasses = ["sentiment-image", "neutral-background"];
        }
      }
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang = "scss" scoped>
h3 {
  color: #ff4301;
}
a {
  float: left;
  color: #bdbdbd;
}

input {
  margin: 5px;
  padding: 8px;
  border-radius: 10px;
}
div {
  text-align: center !important;
}
h5 {
  color: #bdbdbd;
}
body {
  color: #bdbdbd;
  font-family: "Verdana", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
.background {
  background-color: #2c2c2c;
}
.sentiment-image {
  border-radius: 25px;
  margin-top: auto;
  margin-bottom: auto;
}
.negative-background {
  background-color: red;
}
.negative-text {
  color: red;
}

.neutral-background {
  background-color: #bdbdbd;
}
.neutral-text {
  color: #bdbdbd;
}

.positive-background {
  background-color: #42b983;
}
.positive-text {
  color: #42b983;
}

.my-button {
  background-color: #ff4301;
  color: #ffffff;
  border: none;
  text-decoration: none;
  font-family: "Verdana", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  margin: 5px;
  padding: 10px;
  border-radius: 25px;
}
.my-button:hover {
  background-color: #a01d1d;
  cursor: pointer;
}
.sent-text {
  font-size: 9px;
  font-family: "Verdana", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
.subreddit-title {
  font-family: "Verdana", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #ff4301;
}
</style>