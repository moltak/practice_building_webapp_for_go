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


###My go application has been deployed on heroku. 
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

##### Middleware
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


#####HTML Templates
- It has had two parts. Rendering part and template part. 

```go
import (
    "html/template"
    "net/http"
    "path"
)

type Book struct {
    Title  string
    Author string
}

func main() {
    http.HandleFunc("/", ShowBooks)
    http.ListenAndServe(":8080", nil)
}

func ShowBooks(w http.ResponseWriter, r *http.Request) {
    book := Book{"Building Web Apps with Go", "Jeremy Saenz"}

    fp := path.Join("templates", "index.html")
    tmpl, err := template.ParseFiles(fp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := tmpl.Execute(w, book); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

```html
<html>
  <h1>{{ .Title }}</h1>
  <h3>by {{ .Author }}</h3>
</html>
```

######The html/template.ParseFiles has had overhead if you call when every coming request. You can solve it.

#####Renderer
- If I want to rendering JSON and HTML, there is a another way that is using the reneder.
- It has also had two parts. The go file and tmpl file.

```go
import (
	"net/http"

	"gopkg.in/unrolled/render.v1"
)

func main() {
	r := render.New(render.Options{})
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Welcom, visit sub pages now."))
	})

	mux.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
		r.Data(w, http.StatusOK, []byte("Some binary data here."))
	})

	mux.HandleFunc("/json", func(w http.ResponseWriter, req *http.Request) {
		r.JSON(w, http.StatusOK, map[string]string{"hello":"json", "h1":"h2"})
	})

	mux.HandleFunc("/html", func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "example", nil)
	})

	http.ListenAndServe(":3000", mux)
}
```

```html
<h1>Hello {{.}}.</h1>
```