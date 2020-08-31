import Anime from "./components/Anime.vue";
import Downloaded from "./components/Downloaded.vue";
import Settings from "./components/Settings.vue";

export default [
  { path: "/", component: Anime, name: "anime" },
  { path: "/downloaded", component: Downloaded, name: "downloaded" },
  { path: "/settings", component: Settings, name: "settings" }
];
