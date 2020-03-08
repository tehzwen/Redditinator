<template>
  <div>
    <Home v-on:change="handleChange" v-if="currentLocation==='home'" />
    <Data v-on:change="handleChange" v-if="currentLocation==='data'" />
  </div>
</template>

<script>
import Home from "./components/Home";
import Data from "./components/Data";

export default {
  components: {
    Home,
    Data
  },
  data() {
    return {
      topic: null,
      currentLocation: "home"
    };
  },
  methods: {
    testFunc() {
      console.log("here");
    },
    handleChange() {
      if (this.currentLocation === "home") {
        this.currentLocation = "data";
      } else {
        this.currentLocation = "home";
      }
    }
  }
};
</script>

<style lang="scss" scoped>
p {
  font-size: 20px;
}
</style>
