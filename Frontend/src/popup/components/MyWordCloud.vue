<template>
  <div style="width:500px; height:500px;">
    <vue-word-cloud
      :words="this.words"
      :color="['OrangeRed']"
      :font-family="['Verdana']"
      :spacing="0.75"
      :font-size-ratio="25"
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
    return {};
  },
  mounted() {
    let port = chrome.extension.connect({
      name: "Sample Communication"
    });
    this.port = port;
    this.port.postMessage(this.sentiment);
  },
  methods: {}
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
