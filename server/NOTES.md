 ## Notes on Obey

### To get v1 working
 * Config from yaml, env
   * https://dev.to/ilyakaznacheev/a-clean-way-to-pass-configs-in-a-go-application-1g64
   * https://github.com/jinzhu/configor
   * https://github.com/kelseyhightower/envconfig
* Types for Environment, Struct √
* List push, JSON type √
* Iterate over objects √
* Pick http client maybe? √
* http fetch in goroutine √

* work out concurrency model in fetch, ie defer response until all goroutines are finished
   * https://notes.shichao.io/gopl/ch9/

   * Concurrency Notes:
      1. Set up channel for each env
      2. spawn a goroutine for each env
      3. each envs goroutine uses parallel requests.
      4. when requests are done close env channel.
      5. build json out of each env channel data
      https://gobyexample.com/closing-channels

* Reference:
  * https://github.com/gadabout/ohai/blob/master/app/services/services_fetcher.rb
  * https://stackoverflow.com/questions/17539407/how-to-import-local-packages-without-gopath
  * https://medium.com/mindorks/create-projects-independent-of-gopath-using-go-modules-802260cdfb51
  * https://github.com/gin-gonic/gin

# To get v2 working
* ENV vars
* Dockerfile
* websocket streaming
* register routes
    * https://apoorvam.github.io/blog/2017/golang-json-marshal-slice-as-empty-array-not-null/
    * environment √
    * service √
    * https://stackoverflow.com/questions/42967235/golang-gin-gonic-split-routes-into-multiple-files
* redis
  * redis json
    * https://github.com/KromDaniel/rejonson
    * https://www.sohamkamani.com/blog/2017/10/18/parsing-json-in-golang/#parsing-json-strings
    * https://medium.com/@irshadhasmat/*
    * golang-simple-json-parsing-using-empty-interface-and-without-struct-in-go-language-e56d0e69968
    * https://www.restapiexample.com/golang-tutorial/marshal-and-unmarshal-of-struct-data-using-golang/
  *
* polling
  * https://stackoverflow.com/questions/16903348/scheduled-polling-task-in-go
