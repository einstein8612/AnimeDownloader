<template>
  <div id="navigation-wrapper">
    <div class="navigation-item">
      <router-link tag="h3" to="/"> Anime </router-link>
    </div>
    <div class="navigation-item">
      <router-link tag="h3" to="/downloaded"> Downloaded </router-link>
    </div>
    <div class="navigation-item">
      <router-link tag="h3" to="/settings"> Settings </router-link>
    </div>
  </div>
</template>

<script>
export default {
  name: "Navigation"
};
</script>

<style scoped lang="scss">
#navigation-wrapper {
  width: 250px;
  background-color: rgb(47, 49, 54);
}

.navigation-item {
  width: 250px;
  padding: 10px;
  font-size: 25px;
  color: white;
  h3:hover {
    cursor: pointer;
  }
}
</style>
