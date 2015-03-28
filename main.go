package main

import (
	"flag"
	"log"
	"net/http"
	// "os"
	"path/filepath"
	"sync"
	"text/template"
	// "trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// templateHandler: A HTML Template Handler
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP: Handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("ChatAuth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var host = flag.String("host", ":8080", "The host of the application.")
	flag.Parse() // parse the flags

	// Setup OAuth
	gomniauth.SetSecurityKey("test")
	gomniauth.WithProviders(
		facebook.New("FB KEY", "FB SECRET",
			"http://localhost:8080/auth/callback/facebook",
		),
		github.New("GITHUB KEY", "GITHUB SECRET",
			"http://localhost:8080/auth/callback/github",
		),
		google.New("GOOGLE KEY", "GOOGLE SECRET",
			"http://localhost:8080/auth/callback/google",
		),
	)

	r := newRoom()
	// r.tracer = trace.New(os.Stdout)
	// Routes
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets/"))))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// Start the room
	go r.run()
	// start the web server
	log.Println("Starting web server on", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
