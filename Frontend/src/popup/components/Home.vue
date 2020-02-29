<template>
  <div class="background">
    <h3>Redditinator</h3>
    <div>
      <input v-model="topic" />
      <button v-on:click="testFunc">Analyze</button>
    </div>
  </div>
</template>

<script>
export default {
  name: "Home",
  data: () => {
    return {
      topic: null
    };
  },
  mounted() {},
  methods: {
    testFunc() {
      alert(this.topic)
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  color: #FF4301;
}
p {
  color: #42b983;
}
a {
  float: left;
  color: #bdbdbd;
}
button {
  background-color: #FF4301;
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
button:hover {
  background-color: #bdbdbd;
  cursor: pointer;
}

.background {
    background-color: #2c2c2c;
}
</style>