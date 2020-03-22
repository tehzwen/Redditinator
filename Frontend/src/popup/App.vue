<template>
  <div>
    <Home v-on:page="pageChange" v-if="currentLocation==='home'" />
    <Data v-bind:state="dataState" v-on:page="pageChange" v-if="currentLocation==='data'" />
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
      currentLocation: "home",
      dataState: {},
      homeState: {}
    };
  },
  methods: {
    pageChange(val) {
      if (val.page === "data") {
        this.dataState = { subreddit: val.state };
      } else if (val.page === "home") {
        this.homeState = val.state;
      }

      this.currentLocation = val.page;
    }
  }
};
</script>

<style lang="scss">
.row {
  display: flex;
  flex-direction: row;
  flex: 1;
}
.col {
  display: flex;
  flex-direction: column;
  flex: 1;
}
</style>
