<template>
  <div id="app">
    <nav class="navbar is-white">
      <div class="navbar-brand">
        <div class="navbar-item"><div class="block"></div></div>
        <div class="block">
          <img
            src="./assets/obey-logo.svg"
            alt="Peek Obey"
            height="90px"
            width="90px"
          />
        </div>
        <div class="navbar-item">
          <div class="tags has-addons"></div>
        </div>
        <div class="navbar-item"></div>
      </div>
      <div class="navbar-start">
        <div class="navbar-item"></div>
      </div>
      <div class="navbar-end">
        <div class="navbar-item">
          <div class="tags has-addons">
            <span class="tag is-dark is-medium">System Status</span>
            <a class="tag is-warning is-medium is-clickable">
              <b-icon
                :icon="statusIcon"
                size="is-small"
                :type="statusColor"
              ></b-icon>
            </a>
          </div>
        </div>
      </div>
    </nav>
    <!-- <div class="block"></div>     -->
    <!-- <b-navbar>
      <template #brand>
        <b-navbar-item href="#">
          <img src="./assets/obey-logo.png" alt="Peek Obey" />
        </b-navbar-item>
        <b-navbar-item
          ><h1><strong>OBEY</strong></h1>
        </b-navbar-item>
      </template>
      <template #start></template>
      <template #end>
        <b-navbar-item tag="div">
          <div class="buttons is-rounded">
            <a
              class="button is-info is-light is-rounded special-rounded is-noto"
            >
              Sign In
            </a>
          </div>
        </b-navbar-item>
      </template>
    </b-navbar> -->
    <b-tabs v-model="activeTab" position="is-centered">
      <b-tab-item label="Services" icon="chart-scatter-plot-hexbin">
        <div class="container">
          <section class="section">
            <div class="box search-controls">
              <div class="container search-controls">
                <div class="buttons is-right">
                  <b-tooltip
                    position="is-top"
                    type="is-primary is-light"
                    :delay="600"
                    multilined
                  >
                    <div
                      class="button is-warning is-light is-noto"
                      @click="refreshData"
                    >
                      <b-icon
                        icon="refresh"
                        size="is-medium"
                        type="is-dark"
                      ></b-icon>
                    </div>
                    <template v-slot:content>
                      <p><b>Refresh Service Data</b></p>
                      <i
                        >Fetches latest service data from workers if
                        applicable.</i
                      >
                    </template>
                  </b-tooltip>
                  <b-tooltip
                    position="is-top"
                    type="is-primary is-light"
                    :delay="600"
                    multilined
                  >
                    <div
                      class="button is-danger is-light is-noto"
                      @click="clearTableData"
                    >
                      <!-- <span>Clear Table</span> -->
                      <b-icon
                        icon="close"
                        size="is-medium"
                        type="is-dark"
                      ></b-icon>
                    </div>
                    <template v-slot:content>
                      <p><b>Clear Service Data</b></p>
                      <i>Clear table data.</i>.
                    </template>
                  </b-tooltip>
                </div>
                <div class="block"></div>
                <div class="block"></div>
              </div>

              <div class="block">
                <b-field label="Services">
                  <b-taginput
                    v-model="tags"
                    :data="computedTags"
                    autocomplete
                    :allow-duplicates="allowDuplicates"
                    :allow-new="allowNew"
                    :open-on-focus="openOnFocus"
                    icon="chevron-right"
                    type="is-warning is-dark"
                    placeholder="Add Services"
                    @typing="getFilteredTags"
                    @remove="clearTableData"
                  ></b-taginput>
                </b-field>
              </div>
              <div class="block">
                <b-field label="Environments">
                  <b-taginput
                    v-model="envs"
                    :data="computedEnvs"
                    autocomplete
                    :allow-new="allowNew"
                    :open-on-focus="openOnFocus"
                    icon="chevron-right"
                    type="is-success is-dark"
                    placeholder="Add Environments"
                    @typing="getFilteredEnvs"
                    @remove="clearTableData"
                  ></b-taginput>
                </b-field>
              </div>

              <div class="buttons">
                <a
                  class="button is-info is-light is-fullwidth is-noto"
                  @click="buildCompareData2"
                >
                  Compare
                </a>
              </div>
            </div>
          </section>
          <section v-show="!tableIsEmpty" class="section table-holder">
            <table class="table is-striped is-hoverable is-fullwidth">
              <thead>
                <tr>
                  <th>service</th>
                  <th v-for="env in envs" :key="env">
                    {{ env }}
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(object, index) in computedTable" :key="index">
                  <td v-for="(value, i) in object" :key="i">{{ value }}</td>
                </tr>
              </tbody>
            </table>
          </section>
        </div>
      </b-tab-item>
      <b-tab-item label="Diagnostics" icon="chart-line">
        <div class="section">
          <div
            class="container is-flex is-flex-wrap-wrap is-flex-direction-row is-justify-content-center"
          >
            <Worker
              v-for="(worker, index) in registeredWorkers"
              :key="index"
              :worker="worker"
              class="item"
            />
          </div>
        </div>
      </b-tab-item>
      <b-tab-item label="Release Management" icon="chart-gantt"></b-tab-item>
    </b-tabs>
  </div>
</template>

<script>
import { mapState } from "vuex";
import Worker from "./components/diagnostics/Worker.vue";

export default {
  components: {
    Worker,
  },

  data() {
    return {
      filteredTags: [],
      filteredEnvs: [],
      isSelectOnly: false,
      tags: [],
      envs: [],
      allowNew: false,
      allowDuplicates: true,
      openOnFocus: false,
      refreshServices: false,
      isLoadingData: true,
      tableIsEmpty: true,
      tableIsLoading: false,
      tableData: [],
      tableData2: [],
      isBordered: false,
      isStriped: true,
      isHoverable: true,
      isFocusable: true,
      isNarrowed: true,
      activeTab: 0,
      isAuthenticated: false,
      authenticatedUser: "Authenticate",
      statusIcon: "weather-partly-cloudy",
      statusColor: "is-dark",
      registeredWorkers: [
        {
          name: "prod",
          ip: "158.84.106.155",
          type: "k8s",
          uptime: "89:29:17",
          services: 14,
          trend: [1, 12, 9, 13, 10, 8, 4],
        },
        {
          name: "stage",
          ip: "223.103.31.154",
          type: "k8s",
          uptime: "132:13:27",
          services: 16,
          trend: [0, 3, 5, 7, 9, 11, 2],
        },
        {
          name: "qa",
          ip: "112.142.176.156",
          type: "peek-stack",
          uptime: "12:83:00",
          services: 10,
          trend: [3, 2, 8, 2, 8, 1, 10],
        },
        {
          name: "dev",
          ip: "189.228.54.12",
          type: "peek-stack",
          uptime: "02:93:00",
          services: 8,
          trend: [1, 3, 2, 3, 4, 0, 4],
        },
      ],
    };
  },
  methods: {
    getFilteredTags(text) {
      let allAtags = [
        ...new Set(this.serviceData.map((x) => x.services.map((s) => s.name))),
      ].flat();
      // need to dedupe again
      let result = [...new Set(allAtags)].filter(
        (option) =>
          option.toString().toLowerCase().indexOf(text.toLowerCase()) >= 0
      );
      this.filteredTags = result;
      this.$store.commit("updateComputedTags", result);
      console.log("computed tags:", this.computedTags);
    },

    getFilteredEnvs(text) {
      let result = [...new Set(this.serviceData.map((x) => x.name))].filter(
        (option) =>
          option.toString().toLowerCase().indexOf(text.toLowerCase()) >= 0
      );
      this.filteredEnvs = result;
      this.$store.commit("updateComputedEnvs", result);
      console.log("computed envs:", this.computedEnvs);
    },

    buildCompareData2() {
      this.tableIsLoading = true;
      this.tableData2 = [];
      this.tableIsEmpty = true;
      let fe = this.envs
        .map((env) => this.serviceData.filter((option) => option.name == env))
        .flat();
      let composedTable = [];
      this.tags.forEach((tag) => composedTable.push({ service: tag }));
      composedTable.map((row) => {
        fe.map((env) => {
          let filteredServices = env.services.filter(
            (option) => option.name == row.service
          );

          let version;
          if (!filteredServices.length) {
            version = "Not Found";
          } else {
            version = filteredServices[0].version;
          }
          row[env.name] = version;
        });
      });

      this.tableData2 = composedTable.map((row) => Object.values(row));
      this.$store.commit("updateComputedTable", this.tableData2);
      console.log("computedTable", this.computedTable);
      this.tableIsEmpty = false;
      this.tableIsLoading = false;
    },

    clearTableData() {
      this.tableData2 = [];
      this.$store.commit("updateComputedTable", []);
      this.$store.dispatch("clearAllData");
      this.tableIsEmpty = true;
      this.tableIsLoading = false;
    },

    async refreshData() {
      let response = await this.$store.dispatch("callMockApiAsync");
      this.$buefy.snackbar.open({
        duration: 2000,
        message: `Data refreshed- New Env: ${response.name}`,
        type: "is-warning",
        position: "is-bottom-left",
        actionText: "Dismiss",
        queue: false,
      });
    },
    // this is dumb don't use it
    filterArray(array, fields, value) {
      array = array.filter((item) => {
        const found = fields.every((field, index) => {
          return item[field] && item[field] == value[index];
        });
        return found;
      });
      return array;
    },
    // pseudo
    authenticate() {
      this.isAuthenticated = true;
      this.authenticatedUser = "bgroupe";
    },
  },
  mounted() {},
  async created() {
    // let response = await fetch("/api/environments");
    // console.log("async response:", response.json());
    let r = await this.$store.dispatch("callWorkerList");
    console.log(r);
    this.registeredWorkers = this.registeredWorkers.concat(this.workerList);

    console.log(this.registeredWorkers);
  },
  computed: {
    ...mapState([
      "serviceData",
      "computedTable",
      "computedEnvs",
      "computedTags",
      "workerList",
    ]),
  },
};
</script>

<style lang="scss">
$radius-rounded: 290486px !default;
@font-face {
  font-family: "Noto Mono";
  src: url("~@/assets/fonts/NotoMono-Regular.ttf") format("ttf");
}

#app {
  font-family: "Noto Mono";
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #4a4a4a;
}

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }

  .search-controls {
    padding: 20px;
  }

  .container {
    // display: grid;
    // grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    grid-gap: 1.5em;
    // row-gap: 12px;
  }

  .card {
    height: max-content;
  }

  // .button {
  //   &.special-rounded {
  //     border-radius: 290486px;
  //   }
  // }

  .special-rounded {
    border-radius: 290486px;
  }

  .special-logo {
    background: url("./assets/obey.png");
    background-size: cover;
  }

  .is-noto {
    font-family: "Noto Mono";
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

  .logo {
    width: 150px;
    // margin-left: 20px;
    // padding: 30px;
  }

  .has-left-vertical-divider {
    border-left: 6px solid gray;
    height: 36px;
  }

  .test-button {
    border-radius: 75%;
    border: none;
    cursor: pointer;
    color: white;
    background: #444;
  }
}
</style>
