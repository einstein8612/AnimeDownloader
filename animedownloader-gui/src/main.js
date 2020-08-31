import Vue from "vue";
import VueRouter from "vue-router";
import App from "./App.vue";

import routes from "./routes.js";

const router = new VueRouter({
  mode: "history",
  routes
});
Vue.use(VueRouter);

new Vue({
  router,
  render: h => h(App)
}).$mount("#app");
