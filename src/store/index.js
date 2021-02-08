import Vue from "vue";
import Vuex from "vuex";
import "crypto";
import fakerStatic from "faker";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    serviceData: require("@/data/env2.json"),
    computedEnvs: [],
    computedTags: [],
    computedTable: []
  },
  mutations: {
    updateServiceData(state, newData) {
      state.serviceData = newData;
    },

    addMockServiceData(state, newDataPayload) {
      state.serviceData.push(newDataPayload);
    },

    updateComputedTags(state, payload) {
      state.computedTags = payload;
    },
    updateComputedEnvs(state, payload) {
      state.computedEnvs = payload;
    },
    updateComputedTable(state, payload) {
      state.computedTable = payload;
    }
  },
  actions: {
    callMockApi({ commit }) {
      return new Promise((resolve, reject) => {
        let newEnv = generateMockEnv();
        if (newEnv) {
          console.log(newEnv);
          commit("addMockServiceData", newEnv);
          resolve(newEnv);
        } else {
          reject({ error: "failureGenerating mock object" });
        }
      });
    },

    clearAllData(context) {
      context.commit("updateComputedEnvs", []);
      context.commit("updateComputedTags", []);
      context.commit("updateComputedTable", []);
    }
  },
  modules: {}
});

function generateMockEnv() {
  return {
    name: fakerStatic.internet.domainName().toLowerCase(),
    services: [
      {
        name: fakerStatic
          .fake("{{commerce.color}}-{{name.firstName}}")
          .toLowerCase(),
        version: fakerStatic.random.uuid()
      }
    ]
  };
}
