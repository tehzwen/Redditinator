<template>
<div class="background">
  <img src="/assets/redditinator.png" />
  <div>
    <input
      id="subredditSearch"
      placeholder="subreddit"
      v-model="subreddit.name"
      v-on:input="handleSubredditSearch"
    />
    <button v-if="searchMade === true" v-on:click="handleAnalyzeSubmit">Analyze!</button>
    <div v-if="searched.length > 0">
      <div v-for="search in searched" v-bind:key="search.id">
        <span v-on:click="setSearched($event, search)" class="searched-link">{{search.name}}</span>
      </div>
    </div>
  </div>
  <body>Type in a Subreddit and click Analyze!</body>
</div>
</template>

<script>
import Axios from "axios";

export default {
  name: "Home",
  data: () => {
    return {
      subreddit: {},
      port: null,
      searched: [],
      searchMade: false
    };
  },
  mounted() {
    let port = chrome.extension.connect({
      name: "Sample Communication"
    });
    this.port = port;
    document.getElementById("subredditSearch").focus(); //forces the extension to focus the input right away
    // Axios.get("http://167.172.132.5:4000/posts?subreddit=alberta")
    //   .then(res => {
    //     port.postMessage(res);
    //   })
    //   .catch(err => {
    //     port.postMessage(err);
    //   });
  },
  methods: {
    setSearched(e, subreddit) {
      this.subreddit = subreddit;
      this.searched = [];
      this.searchMade = true;
    },
    handleSubredditSearch(e) {
      if (e.target.value !== "") {
        this.searchMade = false;
        Axios.get(
          "http://167.172.132.5:4000/subreddits?subreddit=" + e.target.value
        )
          .then(res => {
            this.searched = res.data;
          })
          .catch(err => {
            this.port.postMessage(err);
          });
      } else {
        this.searched = [];
      }
    },
    handleAnalyzeSubmit() {
      this.$emit("page", { page: "data", state: this.subreddit });
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

.searched-link {
  color: #ff4301;
  font-size: 18px;
}

.searched-link:hover {
  cursor: pointer;
  opacity: 70%;
}
</style>