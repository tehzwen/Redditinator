<template>
  <Home />
</template>

<script>
import Home from "./components/Home";
export default {
  components: {
    Home
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
