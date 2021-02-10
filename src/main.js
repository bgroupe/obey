import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

import Trend from "vuetrend";

import { makeServer } from "./data/mock-server";

import Buefy from "buefy";
import "buefy/dist/buefy.css";

if (process.env.NODE_ENV === "development") {
  makeServer();
}

Vue.use(Buefy);
Vue.use(Trend);

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
