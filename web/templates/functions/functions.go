package functions

import (
	"errors"
	"html/template"
	"net/http"
	"upsilon_cities_go/lib/cities/corporation"
	"upsilon_cities_go/lib/cities/user"
	"upsilon_cities_go/web/tools"
)

//PreLoadFunctions add function at parse time.
//Note If you plan to use custom/local functions, you NEED to add them here first.
func PreLoadFunctions(t *template.Template) {
	fns := make(template.FuncMap)

	fns["IsLogged"] = func() bool { return false }
	fns["IsAdmin"] = func() bool { return false }
	fns["CurrentUser"] = func() (*user.User, error) { return nil, errors.New("not implemented yet") }
	fns["CurrentUserID"] = func() (int, error) { return 0, errors.New("not implemented yet") }
	fns["GetRouter"] = tools.GetRouter
	fns["CurrentCorpID"] = func() (int, error) { return 0, errors.New("not implemented yet") }
	fns["CurrentCorp"] = func() (*corporation.Corporation, error) { return nil, errors.New("not implemented yet") }

	t = t.Funcs(fns)
}

//LoadFunctions add functions to the template
//should find a way to dynamically add functions ...
func LoadFunctions(w http.ResponseWriter, req *http.Request, t *template.Template, fns template.FuncMap) {
	// add generic functions ...

	fns["IsLogged"] = IsLogged(w, req)
	fns["IsAdmin"] = IsAdmin(w, req)
	fns["CurrentUser"] = CurrentUser(w, req)
	fns["CurrentUserID"] = CurrentUser(w, req)
	fns["GetRouter"] = tools.GetRouter
	fns["CurrentCorpID"] = CurrentCorpID(w, req)
	fns["CurrentCorp"] = CurrentCorp(w, req)

	t = t.Funcs(fns)
}

//IsLogged Function generator
func IsLogged(w http.ResponseWriter, req *http.Request) func() bool {
	return func() bool {
		return tools.IsLogged(req)
	}
}

//IsAdmin Function generator
func IsAdmin(w http.ResponseWriter, req *http.Request) func() bool {
	return func() bool {
		return tools.IsAdmin(req)
	}
}

//CurrentUser Function generator
func CurrentUser(w http.ResponseWriter, req *http.Request) func() (*user.User, error) {
	return func() (*user.User, error) {
		return tools.CurrentUser(req)
	}
}

//CurrentUserID Function generator
func CurrentUserID(w http.ResponseWriter, req *http.Request) func() (int, error) {
	return func() (int, error) {
		return tools.CurrentUserID(req)
	}
}

//CurrentCorpID Function generator
func CurrentCorpID(w http.ResponseWriter, req *http.Request) func() (int, error) {
	return func() (int, error) {
		return tools.CurrentCorpID(req)
	}
}

//CurrentCorp Function generator
func CurrentCorp(w http.ResponseWriter, req *http.Request) func() (*corporation.Corporation, error) {
	return func() (*corporation.Corporation, error) {
		return tools.CurrentCorp(req)
	}
}
