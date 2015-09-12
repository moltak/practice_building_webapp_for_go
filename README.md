#Practice building web app for go

###I have been following blow link and started at Sep 12 2015
https://www.gitbook.com/book/codegangsta/building-web-apps-with-go/details

_ _ _

###It has been deployed on heroku. 
1. procFile has been added.
  - touch procFile 
  - copy & paste -> web: BuildingWebApp
2. for dependencies
  - https://github.com/tools/godep
  - godep save -r

### Sub modules
#####Routing
- https://github.com/julienschmidt/httprouter

```go
func main() {
    r := httprouter.New()
    r.GET("/", HomeHandler)
}

// handler
func HomeHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(rw, "Home")
}

```

