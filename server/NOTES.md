 ## Notes on Obey

### To get v1 working
 * Config from yaml, env
   * https://dev.to/ilyakaznacheev/a-clean-way-to-pass-configs-in-a-go-application-1g64
   * https://github.com/jinzhu/configor
   * https://github.com/kelseyhightower/envconfig
* Types for Environment, Struct √
* List push, JSON type √
* Iterate over objects √
* Pick http client maybe?
* http fetch in goroutine

* work out concurrency model in fetch, ie defer response until all goroutines are finished
   * https://notes.shichao.io/gopl/ch9/

* Reference:
  * https://github.com/gadabout/ohai/blob/master/app/services/services_fetcher.rb
  * https://stackoverflow.com/questions/17539407/how-to-import-local-packages-without-gopath
  * https://medium.com/mindorks/create-projects-independent-of-gopath-using-go-modules-802260cdfb51
  * https://github.com/gin-gonic/gin

# To get v2 working
* register routes
* websocket streaming
* redis
