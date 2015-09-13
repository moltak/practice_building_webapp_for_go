#Practice building web app for go

###I have been following blow link and started at Sep 12 2015
https://www.gitbook.com/book/codegangsta/building-web-apps-with-go/details
#####계속 미루고 있었던 Go web app gitbook을 보고 실습. 
#####Go 문법은 https://tour.golang.org/ 에서 미리 조금 공부해놓고 시작했음. (문법 잘 모름....)
#####Go rotine, channel에 대해선 잘 모르지만 gitbook을 따라하는데는 큰 문제가 없었음. 하지만 interface, structure, closure, pointer 는 알고 있어야함. defer도 gitbook 따라하다 알게됐음. (defer 짱!)

_ _ _


###Editor
#####Intellij
- intellij 는 평소에도 많이 사용하는 익숙한 툴이고 go plugin도 잘 되어 있어서 선택했음.
- intellij에서 go plugin을 사용하려면 respository 를 추가해야함.
- $GOPATH를 잘 설정해놓으면 plugin 설치 후 바로 go 코딩을 시작할 수 있음.
1. Add a repository 'https://plugins.jetbrains.com/plugins/alpha/5047'
2. Install Go-lang intellij plugin
  - https://github.com/go-lang-plugin-org/go-lang-idea-plugin
3. Set $GOPATH on intellij preference -> Language & Frameworks -> Go Library -> Global Library


#####중간에 heroku에 deploy 하는 곳이 나옴. 그런데 gitbook이 만들어진지 조금 오래되서 똑같이 따라하면 잘 안됨. 아래 godep를 사용하길 추천함.
1. procFile has been added.
  - touch procFile 
  - copy & paste -> web: BuildingWebApp
  - https://devcenter.heroku.com/articles/getting-started-with-go#introduction
2. for dependencies
  - https://github.com/tools/godep
  - godep save -r

- - -

####Gitbook에 나온 내용들을 다음에 쉽게 보고 따라하려고 정리했음.

- - - -

### Sub modules
#####Routing
- 아주 간단하게 써볼 수 있음 - > 내용 추가해야함.
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
- "github.com/codegangsta/negroni" -> 내용 추가해야함.

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

#####JSON 
- struct와 go 기본 패키지를 이용해서 쉽게 json을 rendering할 수 있음.

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

######Gitbook에 숙제가 있었는데 해결 못했음.... Go 문법 어려워...
- Instead of using the json.Marshal method, try using the json.Encoder API.

#####HTML Templates
- Html에 rendering 하는 부분임.

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

######template.ParseFiles()를 request마다 수행하므로 overhead가 발생함. 저 부분을 해결하라는 숙제가 있었는데, 전역변수로 해결 했음. 더 좋은 방법이 있겠지..
######The html/template.ParseFiles has had overhead if you call when every coming request. You can solve it.


#####Renderer
- JSON이나 HTML을 rendering 할 때, Go 기본 패키지도 좋지만 open source 라이브러리가 많이 나와있음. gitbook에선 "gopkg.in/unrolled/render.v1"만 사용하지만 martini renderer를 많이 사용하는 것 같음.
- tmpl이라는 템플릿 파일을 만들어서 렌더링함.

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
- 테스트 코드의 boiler plate 코드가 엄청 많은것 같음. 조금 더 쉽게 하는 라이브러리가 있을텐데 (찾아보진 않았음.. 아직..)

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

#####negroni 를 테스트하는 부분

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

- 이 부분이 문법으로도 제일 어려웠음. 뭔가 이해가 될듯 말듯 시간이 걸린 부분. type이 func도 될 수 있다는게 엄청 신기했다...
- handler ```
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}``` 를 이용해서 직접 구현했다가 지웠다가 여러가지 트라이해봤음.. 공부는 제일 많이된듯.

- ```type MyController struct {
	AppController
	*render.Render
}``` 이 부분도 엄청 신기한데, field 명을 입력안하고 타입만 쓸 경우 해당 기능을 호출 할 때 자신의 함수인것처럼 하면 되더라.
- 아래 main() 처럼 c.Action을 바로 사용할 수 있음. c.Action은 c.AppController.Action을 줄여서 사용.
```go
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
- 여기도 더 좋은 라이브러를 찾아야할거 같음.. 

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

- - -
#####Go에 대해서 조금 더 익숙해 질 수 있었고 web app을 만드는 것에 대해 살짝.. 발을 담가 볼 수 있었음. 대부분은 재밌었는데 문법을 제대로 몰라서 고생한 부분도 많았고 답답한 부분도 많았음. 