package services

import (
	"net/http"
	"encoding/json"
	"github.com/gadabout/obey/config"
	"fmt"
)

type Service struct {
	Metadata struct {
		Branch string
		Commit string
	}
	Version string
}

type Environment struct {
	Name     string
	Services []Service
}

type Manifest struct {
	CheckEnvironments []CheckEnv
}

type CheckEnv struct {
	Name        string
	ServiceUrls []string
}

func main() {
	// stageUrls := []string{
	// 	"https://noreaga.stage.peek.com/version.json",
	// 	"https://pro.stage.peek.com/version.json",
	// 	"https://notify.stage.peek.com/version.json",
	// 	"https://pro-app.stage.peek.com/version.json",
	// 	"https://book.stage.peek.com/version.json",
	// 	"https://queen.stage.peek.com/version.json",
	// }
	// c := make(chan Service)
	// thing := fetchServices(stageUrls, c)
	manifest := Manifest{
		CheckEnvironments: []CheckEnv{
			{
				Name: "stage",
				ServiceUrls: []string{
					"https://noreaga.stage.peek.com/version.json",
					"https://pro.stage.peek.com/version.json",
					"https://notify.stage.peek.com/version.json",
					"https://pro-app.stage.peek.com/version.json",
					"https://book.stage.peek.com/version.json",
					"https://queen.stage.peek.com/version.json",
				},
			},

			{
				Name: "prod",
				ServiceUrls: []string{
					"https://noreaga.peek.com/version.json",
					"https://pro.peek.com/version.json",
					"https://notify.peek.com/version.json",
					"https://pro-app.peek.com/version.json",
					"https://book.peek.com/version.json",
					"https://queen.peek.com/version.json",
				},
			},
		},
	}

	ec := make(chan Environment)
	fetchedEnvironments := []Environment{}
	for _, checkEnv := range manifest.CheckEnvironments {
		fmt.Println("fetching environment", checkEnv.Name)
		go fetchEnvironment(checkEnv, ec)
		fetchedEnvironments = append(fetchedEnvironments, <-ec)
	}
	// ^^^
	// just do this but with envs now
	// wrap this block in an env block, wrap in a func

	fmt.Println(fetchedEnvironments)
}

func fetchEnvironment(checkEnv CheckEnv, ec chan Environment) {
	c := make(chan Service)
	fetchedServices := fetchServices(checkEnv.ServiceUrls, c)
	e := Environment{
		Name:     checkEnv.Name,
		Services: fetchedServices,
	}
	ec <- e
}

func fetchServices(serviceUrls []string, c chan Service) []Service {
	fetchedServices := []Service{}

	for _, url := range serviceUrls {
		fmt.Println("checking host", url)
		go getServiceVersion(url, c)
		fetchedServices = append(fetchedServices, <-c)
	}
	return fetchedServices
}

func getServiceVersion(url string, c chan Service) {
	fmt.Println("entered")

	res, err := http.Get(url)

	if err != nil {
		panic(err)
	} else {
		fmt.Println(url, "checked")
	}

	defer res.Body.Close()
	var svc Service
	json.NewDecoder(res.Body).Decode(&svc)
	c <- svc

}