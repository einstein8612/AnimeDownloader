<template>
  <div id="app">
    <Navigation />
    <router-view></router-view>
  </div>
</template>

<script>
import Navigation from "./components/Navigation.vue";

export default {
  name: "App",
  components: {
    Navigation
  }
};
</script>

<style>
* {
  font-family: Consolas, Helvetica, Arial, sans-serif;
  margin: 0;
  padding: 0;
  color: rgb(142, 146, 151);
}
#app {
  display: flex;
}
</style>
