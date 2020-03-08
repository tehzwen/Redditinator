<template>
  <div>
    <Home v-on:page="pageChange" v-if="currentLocation==='home'" />
    <Data v-on:page="pageChange" v-if="currentLocation==='data'" />
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
    pageChange(val) {
      this.currentLocation = val.page;
    }
  }
};
</script>

<style lang="scss" scoped>
p {
  font-size: 20px;
}
</style>
