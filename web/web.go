package web

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/song940/acme-go/acme"
)

// go:embed templates/*.html
var Files embed.FS

type H map[string]interface{}

type Server struct {
	acme *acme.Client
}

func NewServer() (server *Server) {
	server = &Server{
		acme: acme.NewClient(nil),
	}
	return
}

// Render renders an HTML template with the provided data.
func (server *Server) Render(w http.ResponseWriter, templateName string, data H) {
	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/"+templateName+".html")
	// Parse templates from embedded file system
	// tmpl, err := template.New("").ParseFS(Files, "layout.html", templateName+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Execute "index.html" within the layout and write to response
	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *Server) Error(w http.ResponseWriter, err error) {
	server.Render(w, "error", H{
		"error": err.Error(),
	})
}

func (s *Server) IndexView(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		s.Render(w, "index", H{})
		return
	}

	r.ParseForm()
	action := r.URL.Query().Get("action")

	if action == "key" {
		key := r.Form.Get("key")
		s.acme.ImportKey(key)
	}

	if action == "directory" {
		s.acme.Config.DirectoryURL = r.FormValue("directory_url")
		s.acme.Directory, _ = s.acme.GetDirectory()
	}

	if action == "login" {
		account_url := r.FormValue("account_url")
		s.acme.AccountURL = account_url
		resp, err := s.acme.GetAccount(account_url)
		log.Println(resp, err)
	}

	if action == "register" {
		email := r.Form.Get("email")
		accountUrl, resp, err := s.acme.Register(&acme.AccountRequest{
			TermsOfServiceAgreed: true,
			Contact: []string{
				"mailto:" + email,
			},
		})
		log.Println(accountUrl, resp, err)
	}

	if action == "order" {
		domain := r.Form.Get("domain")
		orderUrl, resp, err := s.acme.CreateOrder(&acme.OrderRequest{
			Identifiers: []acme.Identifier{
				{
					Type:  "dns",
					Value: domain,
				},
			},
		})
		log.Println(orderUrl, resp, err)
	}

	s.Render(w, "index", H{})
}

func (s *Server) OrderView(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	order, err := s.acme.GetOrder(url)
	if err != nil {
		s.Error(w, err)
		return
	}
	s.Render(w, "order", H{
		"order": order,
	})
}

func (s *Server) AuthzView(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	authz, err := s.acme.GetAuthorization(url)
	if err != nil {
		s.Error(w, err)
		return
	}
	s.Render(w, "authz", H{
		"authz": authz,
	})
}

func (s *Server) Start() {
	router := http.NewServeMux()
	router.HandleFunc("/", s.IndexView)
	router.HandleFunc("/order", s.OrderView)
	router.HandleFunc("/authz", s.AuthzView)
	http.ListenAndServe(":8080", router)
}
