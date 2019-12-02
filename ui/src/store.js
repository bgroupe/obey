import Vue from "vue";
import Vuex from "vuex";
import { loadConfig, parseConfigurationEndpoints }  from './helpers/loadConfig'

Vue.use(Vuex);

const envConfig = loadConfig()
const endpoints = parseConfigurationEndpoints(envConfig)

export default new Vuex.Store({
  state: {
    services: envConfig.services
    environments:
  },
  mutations: {},
  actions: {}
});
