package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type templateData struct {
	StringMap            map[string]string
	IntMap               map[string]int
	FloatMap             map[string]float32
	Data                 map[string]interface{}
	CSRFToken            string
	Flash                string
	Warning              string
	Error                string
	IsAuthenticated      int
	UserID               int
	API                  string
	CssVersion           string
	StripeSecrectKey     string
	StripePublishableKey string
}

var functions = template.FuncMap{
	"formatCurrency": formatCurrency,
}

func formatCurrency(n int) string {
	f := float32(n) / float32(100)
	return fmt.Sprintf("$%.2f", f)
}

//go:embed templates
var templateFS embed.FS

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	td.API = app.config.api
	td.StripeSecrectKey = app.config.stripe.secret
	td.StripePublishableKey = app.config.stripe.key

	if app.Session.Exists(r.Context(), "userID") {
		td.IsAuthenticated = 1
		td.UserID = app.Session.GetInt(r.Context(), "userID")
	} else {
		td.IsAuthenticated = 0
		td.UserID = 0
	}

	return td
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData, partials ...string) error {
	var t *template.Template
	var err error

	templateToRender := fmt.Sprintf("templates/%s.page.html", page)

	_, templateInMap := app.templateCache[templateToRender]

	if app.config.env == "production" && templateInMap {
		t = app.templateCache[templateToRender]
	} else {
		t, err = app.parseTemplate(partials, page, templateToRender)
		if err != nil {
			app.errorLog.Println(err)
			return err
		}
	}

	if td == nil {
		td = &templateData{}
	}

	println(td.API)
	td = app.addDefaultData(td, r)
	err = t.Execute(w, td)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	return nil
}

func (app *application) parseTemplate(partials []string, page, templateToRender string) (*template.Template, error) {
	var t *template.Template
	var err error

	//build templates
	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("templates/%s.partial.html", x)
		}
	}
	if len(partials) > 0 {
		t, err = template.New(fmt.Sprintf("%s.page.html", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.html", strings.Join(partials, ","), templateToRender)
	} else {
		t, err = template.New(fmt.Sprintf("%s.page.html", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.html", templateToRender)
	}

	if err != nil {
		app.errorLog.Println(err)
		return nil, err
	}

	app.templateCache[templateToRender] = t
	return t, nil
}
