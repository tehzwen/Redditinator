<template>
  <div
    style="width:250px; height:150px; margin-right:15px; margin-left:15px; margin-bottom:15px"
  >
    <vue-word-cloud
      :words="this.words"
      :color="[getColor()]"
      :font-family="['Verdana']"
      :spacing="0.75"
      :font-size-ratio="25"
      :animation-duration="250"
    >
      <template slot-scope="{ text, weight, word }">
        <div
          :title="weight"
          style="cursor: pointer;"
          v-on:click="onWordClick(word)"
        >
          {{ text }}
        </div>
      </template>
    </vue-word-cloud>
  </div>
</template>

<script>
import VueWordCloud from "vuewordcloud";

export default {
  props: {
    words: Array,
    sentiment: Object
  },
  components: {
    [VueWordCloud.name]: VueWordCloud
  },
  data: () => {
    return {
      sentimentType: "post"
    };
  },
  mounted() {
    let port = chrome.extension.connect({
      name: "Sample Communication"
    });
    this.port = port;
    //this.port.postMessage(this.sentiment.postSent);
  },
  methods: {
    getColor() {
      //this.port.postMessage(this.sentimentType);
      var currSent =
        this.sentimentType == "post"
          ? this.sentiment.postSent
          : this.sentiment.commentSent;
      return currSent > 0 ? "#42b983" : currSent < 0 ? "#FF0000" : "#bdbdbd";
    },
    changeSentiment(type) {
      this.sentimentType = type;
      this.VueWordCloud.color = getColor();
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
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
  font-family: "Verdana";
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
