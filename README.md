#Practice building web app for go

###I have been following blow link and started at Sep 12 2015
https://www.gitbook.com/book/codegangsta/building-web-apps-with-go/details

_ _ _


###Editor
######I choose the intellij and Go plugins. Intellij is easy to use and convenience.
1. Add a repository 'https://plugins.jetbrains.com/plugins/alpha/5047'
2. Install Go-lang intellij plugin
  - https://github.com/go-lang-plugin-org/go-lang-idea-plugin
3. Set $GOPATH on intellij preference -> Language & Frameworks -> Go Library -> Global Library


###My Go app has been deployed on heroku. 
1. procFile has been added.
  - touch procFile 
  - copy & paste -> web: BuildingWebApp
  - https://devcenter.heroku.com/articles/getting-started-with-go#introduction
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

##### middleware
- "github.com/codegangsta/negroni"

```go
func main() {
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(MyMiddleware),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("public")),
	)
	n.Run(":3000")
}

func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("Loggin on the way there...")

	if r.URL.Query().Get("password") == "secret123" {
		next(rw, r)
	} else {
		http.Error(rw, "Not Authorized", 401)
	}

	log.Println("Loggin on the way back...")
}
```

##### json
- I have to memory the how to make a json struct and json.Marshal

```go
import (
	"encoding/json"
	"net/http"
)

type Book struct {
	Title string `json:"title"` // There is a strange syntax.
	Author string `json:"author"`
}

func ShowBooks(w http.ResponseWriter, r *http.Request) {
	book := Book{"Building Web Apps with Go", "moltak"}

	js, err := json.Marshal(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(js)
}
```

######I can't solve the this quest. -> Instead of using the json.Marshal method, try using the json.Encoder API.