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

#####Test
- Test needs to many codes. I have to find out better way. It is so difficult I commit that test codes to memory.

```go
import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloWorld)
	http.ListenAndServe(":3000", nil)
}

func HelloWorld(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		fmt.Fprint(res, "Hello World")
	} else {
		fmt.Fprint(res, "Hello Post")
	}
}
```
```go
package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_App(t *testing.T) {
	ts := httptest.NewServer(App())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	exp := "Before...Hello World...After"

	if exp != string(body) {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}

func Test_AppPost(t *testing.T) {
	ts := httptest.NewServer(App())
	defer ts.Close()

	bodyType := "string"
	res, _ := http.Post(ts.URL, bodyType, nil)
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	exp := "Before...Hello Post...After"
	if exp != string(body) {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}
```

- This is using negroni.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

func main() {
	http.ListenAndServe(":3000", App())
}

func App() http.Handler {
	n := negroni.Classic()

	m := func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		fmt.Fprint(res, "Before...")
		next(res, req)
		fmt.Fprint(res, "...After")
	}
	n.Use(negroni.HandlerFunc(m))

	r := httprouter.New()

	r.GET("/", HelloWorld)
	r.POST("/", HelloPost)
	n.UseHandler(r)
	return n
}

func HelloWorld(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	fmt.Fprint(res, "Hello World")
}

func HelloPost(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	fmt.Fprint(res, "Hello Post")
}
```
```go
package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_App(t *testing.T) {
	ts := httptest.NewServer(App())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	exp := "Before...Hello World...After"

	if exp != string(body) {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}

func Test_AppPost(t *testing.T) {
	ts := httptest.NewServer(App())
	defer ts.Close()

	bodyType := "string"
	res, _ := http.Post(ts.URL, bodyType, nil)
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	exp := "Before...Hello Post...After"
	if exp != string(body) {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}
```

#####Controller

- I've made one controller, controller can process GET, POST.

```go
package main

import (
	"net/http"

	"gopkg.in/unrolled/render.v1"
)


type Action func(res http.ResponseWriter, req *http.Request) error

type AppController struct {}

func (c *AppController) Action(a Action) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if err := a(res, req); err != nil {
			http.Error(res, err.Error(), 500)
		}
	})
}

type MyController struct {
	AppController
	*render.Render
}

func (c *MyController) Index(res http.ResponseWriter, req *http.Request) error {
	if req.Method == "GET" {
		c.JSON(res, 200, map[string]string{"Hello":"GET"})
	} else {
		c.JSON(res, 200, map[string]string{"Hello":"POST"})
	}
	return nil
}

func main() {
	c := &MyController{Render: render.New(render.Options{})}
	http.Handle("/", c.Action(c.Index))
	http.ListenAndServe(":3000", nil)
}

```
```go
package main

import (
	"net/http"
	"testing"
	"net/http/httptest"
	"io/ioutil"

	"gopkg.in/unrolled/render.v1"
)

func Test_ControllerMainPost(t *testing.T) {
	c := &MyController{Render: render.New(render.Options{})}
	ts := httptest.NewServer(c.Action(c.Index))
	defer ts.Close()

	res, _ := http.Post(ts.URL, "string", nil)
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	exp := "{\"Hello\":\"POST\"}"
	if exp != string(body) {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}

func Test_ControllerMainGet(t *testing.T) {
	c := &MyController{Render: render.New(render.Options{})}
	ts := httptest.NewServer(c.Action(c.Index))
	defer ts.Close()

	res, _ := http.Get(ts.URL)
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	exp := "{\"Hello\":\"GET\"}"
	if exp != string(body) {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}
```

#####Database
- I've added database codes.

```go
package main
import (
	"database/sql"
	"log"
	"net/http"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := NewDB()
	log.Println("Listeneing on :3000")
	http.ListenAndServe(":3000", ShowBooks(db))
}

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "example.sqlite")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table if not exists books(title text, author text)")
	if err != nil {
		panic(err)
	}

	return db
}

func ShowBooks(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var title, author string
		err := db.QueryRow("select title, author from books").Scan(&title, &author)
		if err != nil {
			InsertBooks(db)
			panic(err)
		}

		fmt.Fprintf(rw, "The first book is '%s' by '%s'", title, author)
	})
}

func InsertBooks(db *sql.DB) {
	for i := 1; i < 3; i ++ {
		sql := fmt.Sprintf("insert into books values ('book%d', 'author%d')", i, i)
		fmt.Print(sql)
		_, err := db.Exec(sql)
		if err != nil {
			panic(err)
		}
	}
}
```