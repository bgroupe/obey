import Vue from "vue";
import Vuex from "vuex";
import fakerStatic from "faker";
import axios from "axios";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    serviceData: require("@/data/env2.json"),
    computedEnvs: [],
    computedTags: [],
    computedTable: [],
    workerList: []
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
    },
    clearWorkers(state) {
      state.workerList = []
    },

    appendWorker(state, payload) {
      state.workerList.push(payload)
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

    async callMockApiAsync(context) {
      let result = await (await fetch("/api/environments")).json();

      return new Promise((resolve, reject) => {
        let newEnv = result;
        if (newEnv) {
          context.commit("addMockServiceData", newEnv);
          resolve(newEnv);
        } else {
          reject({ error: "failureGenerating mock object" });
        }
      });
    },

    async callWorkerList(context) {
      // POC
      let result = await axios.get("http://127.0.0.1:3000/list")
      return new Promise((resolve, reject) => {
        let workerList = result.data;
        context.commit('clearWorkers')
        if (workerList) {
          Object.keys(workerList).map(key => {
            let newWorker = {
              name: workerList[key].env,
              ip: workerList[key].address,
              type: workerList[key]["env-type"],
              id: workerList[key].id,
              uptime: workerList[key]["launch-time"],
              services: 0,
              trend: [0, 1, 2, 3, 4, 5, 6]
            }

            context.commit('appendWorker', newWorker)
          })
          console.log("newWorkers:", context.state.workerList)
          resolve(workerList)
        } else {
          reject({ error: "failure generating" })
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
