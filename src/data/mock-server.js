import { Server } from "miragejs";
import fakerStatic from "faker";

export function makeServer({ environment = "development" } = {}) {
  let server = new Server({
    environment,

    routes() {
      this.namespace = "api";
      this.get("/environments", () => {
        return generateMockEnv();
      });

      this.get("/workers", () => {
        return {};
      });
    }
  });

  return server;
}

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
