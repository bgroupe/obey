yaml = require("js-yaml");
fs = require("fs");
path = require("path");

const NOT_APPLICABLE = 'N/A Not Scraping'

module.exports =  {
  function loadConfig() {
    let filePath = path.join(__dirname, "..", "..", "..", "environment.yaml");
    let configFile = fs.readFileSync(filePath, "utf8");

    try {
      return yaml.safeLoad(configFile);
    } catch (e) {
      console.log("Could not load config YAML", e);
    }
  },

  function parseServiceEndpoints(config) {
    let envs = {}
    config.environments.forEach(environment => {
        let endpoints = {}
        config.services.forEach(service => {
        let serviceName = service.name
          endpoints[serviceName] = `${serviceName}.${environment.host}/version.json`
        })
        envs[environment.name] = endpoints
    })
    return envs
  }
}
