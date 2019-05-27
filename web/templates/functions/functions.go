package functions

import (
	"errors"
	"html/template"
	"net/http"
	"upsilon_cities_go/lib/cities/user"
	"upsilon_cities_go/web/tools"
)

//PreLoadFunctions add function at parse time.
//Note If you plan to use custom/local functions, you NEED to add them here first.
func PreLoadFunctions(t *template.Template) {
	fns := make(template.FuncMap)

	fns["IsLogged"] = func() bool { return false }
	fns["IsMap"] = func() bool { return false }
	fns["CurrentCorpID"] = func() (int, error) { return 0, errors.New("not implemented yet") }
	fns["CurrentCorpName"] = func() (string, error) { return "", errors.New("not implemented yet") }
	fns["IsAdmin"] = func() bool { return false }
	fns["CurrentUser"] = func() (*user.User, error) { return nil, errors.New("not implemented yet") }
	fns["CurrentUserID"] = func() (int, error) { return 0, errors.New("not implemented yet") }
	fns["GetRouter"] = tools.GetRouter
	fns["CurrentCorpID"] = func() int { return 0 }
	fns["CurrentCorpName"] = func() (string, error) { return "", errors.New("not implemented yet") }

	fns["ErrorAlerts"] = func() string { return "" }
	fns["InfoAlerts"] = func() string { return "" }
	fns["WarningAlerts"] = func() string { return "" }

	t = t.Funcs(fns)
}

//LoadFunctions add functions to the template
//should find a way to dynamically add functions ...
func LoadFunctions(w http.ResponseWriter, req *http.Request, t *template.Template, fns template.FuncMap) {
	// add generic functions ...

	fns["IsLogged"] = IsLogged(w, req)
	fns["IsAdmin"] = IsAdmin(w, req)
	fns["CurrentCorpID"] = CurrentCorpID(w, req)
	fns["IsMap"] = IsMap(w, req)
	fns["CurrentUser"] = CurrentUser(w, req)
	fns["CurrentUserID"] = CurrentUser(w, req)
	fns["CurrentCorpName"] = CurrentCorpName(w, req)
	fns["GetRouter"] = tools.GetRouter
	fns["ErrorAlerts"] = ErrorAlerts(w, req)
	fns["InfoAlerts"] = InfoAlerts(w, req)
	fns["WarningAlerts"] = WarningAlerts(w, req)

	t = t.Funcs(fns)
}

//IsLogged Function generator
func IsLogged(w http.ResponseWriter, req *http.Request) func() bool {
	return func() bool {
		return tools.IsLogged(req)
	}
}

//IsMap Function generator
func IsMap(w http.ResponseWriter, req *http.Request) func() bool {
	return func() bool {
		return tools.IsMap(req)
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
func CurrentCorpID(w http.ResponseWriter, req *http.Request) func() int {
	return func() int {
		cid, _ := tools.CurrentCorpID(req)
		return cid
	}
}

//CurrentCorpName Function generator
func CurrentCorpName(w http.ResponseWriter, req *http.Request) func() (string, error) {
	return func() (string, error) {
		return tools.CurrentCorpName(req)
	}
}

//ErrorAlerts tell whether alerts marked as errors are available.
func ErrorAlerts(w http.ResponseWriter, req *http.Request) func() string {
	return func() string {
		return tools.ErrorAlerts(req)
	}
}

//InfoAlerts tell whether alerts marked as errors are available.
func InfoAlerts(w http.ResponseWriter, req *http.Request) func() string {
	return func() string {
		return tools.InfoAlerts(req)
	}
}

//WarningAlerts tell whether alerts marked as errors are available.
func WarningAlerts(w http.ResponseWriter, req *http.Request) func() string {
	return func() string {
		return tools.WarningAlerts(req)
	}
}
