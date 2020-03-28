<template>
  <div style="width:500px; height:500px;">
    <vue-word-cloud
      :words="[['romance', 19], ['horror', 3], ['fantasy', 7], ['adventure', 3]]"
      :color="([, weight]) => weight > 10 ? 'DeepPink' : weight > 5 ? 'RoyalBlue' : 'Indigo'"
      font-family="Roboto"
      height="500px"
      width="500px"
    />
  </div>
</template>

<script>
import VueWordCloud from "vuewordcloud";

export default {
  props: {},
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
    this.port.postMessage(VueWordCloud);
  },
  methods: {}
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
</style>