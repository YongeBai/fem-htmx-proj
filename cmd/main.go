package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"	
)

type Templates struct{
	templates *template.Template	
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates {
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Contact struct {
	Name string
	Email string
}

func CreateContact(name, email string) Contact {
	return Contact{
		Name: name,
		Email: email,
	}
}

func (d *Data) hasEmail(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

type Contacts = []Contact

type Data struct {
	Contacts Contacts
}

func CreateData() Data {
	return Data{
		Contacts: Contacts{
			CreateContact("John Doe", "jd@gmail.com"),
			CreateContact("Kanye West", "kw@gmail.com"),
		},
	}
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func CreateFormData() FormData {
	return FormData {
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	FormData FormData
}

func CreatePage() Page {
	return Page {
		Data: CreateData(),
		FormData: CreateFormData(),
	}
}

func main() {
	e := echo.New()	
	e.Use(middleware.Logger())

	page := CreatePage()
	e.Renderer = NewTemplates()

	e.GET("/", func(c echo.Context) error {		
		return c.Render(100, "index", page)
	})


	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.Data.hasEmail(email) {
			formData := CreateFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "Email already exists"

			return c.Render(422, "create-contact", formData)
		}

		page.Data.Contacts = append(page.Data.Contacts, CreateContact(name, email))
		return c.Render(200, "contact-list", page.Data)
	})


	e.Logger.Fatal(e.Start(":8080"))
}