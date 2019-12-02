const feathers = require('@feathersjs/feathers')
const app = feathers()

class HelloService {
  constructor() {
      this.names = []
  }

  create(data,params){
      const name = data.name
      this.names.push(name)
      return Promise.resolve(this.helloTo(name))
  }

  find(params){
      return Promise.resolve(this.names.map(this.helloTo))
  }

  helloTo(name){
      return `Hello ${name}!`
  }
}

app.use('hello', new HelloService() )

const helloService = app.service('hello')
const readline = require('readline').createInterface(process.stdin, process.stdout)

readline.setPrompt("what's ur name?: ")
readline.prompt()

readline.on('line', name => {
    const response = name =='everyone' ?
      helloService.find().then(hellos => hellos.forEach(hello => console.log(hello))) :
      helloService.create({ name }).then(hello => console.log(hello))
    response.then(() => readline.prompt())
})
