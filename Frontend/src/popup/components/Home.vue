<template>
<div class="background">
  <img src="/assets/redditinator.png" />
  <div>
    <input v-model="topic" />
    <button v-on:click="testFunc">Analyze!</button>
  </div>
  <body>Type in a Subreddit and click Analyze!</body>
  <button v-on:click="goToData">Go to data!</button>
  <div v-if="sentimentValue !== null">
    <p>{{sentimentValue.SentimentOverall}}</p>
  </div>
</div>
</template>

<script>
import Axios from "axios";

export default {
  name: "Home",
  data: () => {
    return {
      topic: null,
      port: null,
      sentimentValue: null
    };
  },
  mounted() {
    let port = chrome.extension.connect({
      name: "Sample Communication"
    });

    this.port = port;
    Axios.get("http://167.172.132.5:4000/posts?subreddit=alberta")
      .then(res => {
        port.postMessage(res);
      })
      .catch(err => {
        port.postMessage(err);
      });
  },
  methods: {
    testFunc() {
      Axios.post("http://167.172.132.5:4000/sentiment", {
        text: this.topic
      })
        .then(res => {
          this.sentimentValue = res.data;
        })
        .catch(err => {
          this.port.postMessage(err);
        });
    },
    goToData() {
      this.$emit("page", { page: "data" });
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang = "scss" scoped>
h3 {
  color: #ff4301;
}
p {
  color: #42b983;
}
a {
  float: left;
  color: #bdbdbd;
}
button {
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
button:hover {
  background-color: #a01d1d;
  cursor: pointer;
}
.background {
  background-color: #2c2c2c;
}
</style>