package api_test

import (
	"fmt"

	apitrace "github.com/lonnblad/dd-trace-go/contrib/google.golang.org/api"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func Example() {
	// create an oauth2 client suitable for use with the google APIs
	client, _ := apitrace.NewClient(
		// set scopes like this, which will vary depending on the service
		apitrace.WithScopes(cloudresourcemanager.CloudPlatformScope))
	svc, _ := cloudresourcemanager.New(client)

	// call google api methods as usual
	res, _ := svc.Projects.List().Do()
	for _, project := range res.Projects {
		fmt.Println(project.Name)
	}
}
