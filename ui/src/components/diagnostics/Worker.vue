<template>
  <div>
    <div class="box worker">
      <div class="level">
        <div class="level-left">
          <div class="level-item">
            <b-icon
              type="is-primary"
              icon="cube-scan"
              size="is-medium"
            ></b-icon>
          </div>
          <div class="level-item">
            <h2 class="underline is-size-4">{{ worker.name }}</h2>
          </div>
        </div>
      </div>
      <div class="level">
        <div class="level-left">
          <div class="level-item">
            <b-taglist attached>
              <b-tag type="is-warning">
                <h2>Services</h2>
              </b-tag>
              <b-tag type="is-light">
                <h2>
                  <strong>{{ worker.services }}</strong>
                </h2>
              </b-tag>
            </b-taglist>
          </div>
          <div class="level-item">
            <b-taglist attached>
              <b-tag type="is-primary is-light">
                <h2>Status</h2>
              </b-tag>
              <b-tag type="is-success">
                <h2>
                  <strong class="has-text-white">Healthy</strong>
                </h2>
              </b-tag>
            </b-taglist>
          </div>
        </div>
      </div>
      <div class="level">
        <div class="level-left">
          <div class="level-item">
            <b-taglist attached>
              <b-tag type="is-info is-light">
                <b-icon
                  type="is-info"
                  icon="vector-intersection"
                  size="is-small"
                ></b-icon>
              </b-tag>
              <b-tag type="is-light ">
                <h2>
                  <strong>{{ worker.ip }}</strong>
                </h2>
              </b-tag>
            </b-taglist>
          </div>
          <div class="level-item">
            <b-taglist attached>
              <b-tag type="is-success is-light ">
                <b-icon
                  type="is-success"
                  icon="shape-plus"
                  size="is-small"
                ></b-icon>
              </b-tag>
              <b-tag type="is-light ">
                <h2>
                  <strong>{{ worker.type }}</strong>
                </h2>
              </b-tag>
            </b-taglist>
          </div>
        </div>
      </div>
      <div class="level">
        <div class="level-left">
          <div class="level-item">
            <b-taglist attached>
              <b-tag type="is-danger is-light ">
                <b-icon
                  type="is-warning"
                  icon="clock-time-seven"
                  size="is-small"
                ></b-icon>
              </b-tag>
              <b-tag type="is-light  ">
                <h2>
                  <strong>{{ worker.uptime }}</strong>
                </h2>
              </b-tag>
            </b-taglist>
          </div>
        </div>
      </div>
      <div class="level">
        <div class="level-left">
          <div class="level-item">
            <h2>7 Day Trend</h2>
          </div>
        </div>
      </div>

      <trend
        :data="worker.trend"
        :gradient="['#6fa8dc', '#42b983', '#2c3e50']"
        auto-draw
        smooth
      ></trend>
    </div>
  </div>
</template>
<script>
export default {
  name: "Worker",
  props: {
    worker: Object
  }
};
</script>
<style lang="scss">
.worker {
  width: 350px;
  // height: 300px;
  padding: 50px;
  margin: 24px;
  cursor: pointer;
}

.underline {
  text-decoration: underline;
}
</style>
