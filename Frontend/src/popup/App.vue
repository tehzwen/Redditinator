<template>
  <Home />
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
      topic: null
    };
  },
  methods: {
    testFunc() {
      console.log("here");
    }
  }
};
</script>

<style lang="scss" scoped>
p {
  font-size: 20px;
}
</style>
