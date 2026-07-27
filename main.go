// Aegis demo: Gin + cookie session. Target for scanner pipelines.
package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var (
	users = map[string]string{"demo@example.com": "demo1234"}
	items = []item{
		{1, "Widget", 19.99},
		{2, "Gadget", 24.50},
		{3, "Sprocket", 8.75},
	}
)

const loginTmpl = `<!doctype html><title>Login</title>
<h1>Sign in</h1>
{{ if .Error }}<p style="color:red">{{ .Error }}</p>{{ end }}
<form method="post" action="/login">
  <label>Email <input name="email" type="email" required></label><br>
  <label>Password <input name="password" type="password" required></label><br>
  <button type="submit">Sign in</button>
</form>`

const dashTmpl = `<!doctype html><title>Dashboard</title>
<h1>Welcome, {{ .Email }}</h1>
<ul>{{ range .Items }}<li>{{ .Name }} — ${{ .Price }}</li>{{ end }}</ul>
<form method="post" action="/logout"><button>Sign out</button></form>`

func secret() []byte {
	if s := os.Getenv("SESSION_SECRET"); s != "" {
		return []byte(s)
	}
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return []byte(hex.EncodeToString(b))
}

func main() {
	r := gin.Default()

	tmpl, err := loadTemplates()
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(tmpl)

	store := cookie.NewStore(secret())
	r.Use(sessions.Sessions("aegis-demo-go", store))

	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8",
			[]byte("<h1>Aegis demo (Go)</h1><a href='/login'>Sign in</a>"))
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login", gin.H{"Error": ""})
	})

	r.POST("/login", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		if pw, ok := users[email]; ok && pw == password {
			s := sessions.Default(c)
			s.Set("email", email)
			_ = s.Save()
			c.Redirect(http.StatusFound, "/dashboard")
			return
		}
		c.HTML(http.StatusUnauthorized, "login", gin.H{"Error": "Invalid credentials"})
	})

	authed := r.Group("/", func(c *gin.Context) {
		s := sessions.Default(c)
		if s.Get("email") == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
		}
	})

	authed.GET("/dashboard", func(c *gin.Context) {
		s := sessions.Default(c)
		c.HTML(http.StatusOK, "dash", gin.H{"Email": s.Get("email"), "Items": items})
	})

	authed.GET("/api/items", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	r.POST("/logout", func(c *gin.Context) {
		s := sessions.Default(c)
		s.Clear()
		_ = s.Save()
		c.Redirect(http.StatusFound, "/")
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	_ = r.Run(":" + port)
}
