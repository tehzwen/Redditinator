<template>
<div class="background">
  <img src="/assets/redditinator.png" />
  <body>Data Page</body>
  <button v-on:click="goToHome">Go to home!</button>
</div>
</template>

<script>
export default {
  methods: {
    goToHome() {
      this.$emit("change");
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