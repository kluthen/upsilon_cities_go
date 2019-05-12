package templates

// courtesy to https://hackernoon.com/golang-template-2-template-composition-and-how-to-organize-template-files-4cb40bcdf8f6

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"upsilon_cities_go/config"
	"upsilon_cities_go/web/tools"

	"github.com/oxtoacart/bpool"
)

var mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}`

var templateConfig struct {
	TemplateLayoutPath  string
	TemplateIncludePath string
	TemplateSharedPath  string
}

type templateInfo struct {
	tmpl       *template.Template
	baseTmpl   *template.Template
	path       string
	base       string
	lastUpdate time.Time
}

var templates map[string]templateInfo

type sharedInfo struct {
	path       string
	lastUpdate time.Time
}

var layouts map[string]map[string]sharedInfo
var sharedCheck map[string]sharedInfo
var shared []string
var bufpool *bpool.BufferPool

func loadConfiguration() {
	templateConfig.TemplateLayoutPath = config.WEB_LAYOUTS
	templateConfig.TemplateSharedPath = config.WEB_SHARED
	templateConfig.TemplateIncludePath = config.WEB_TEMPLATES
}

func paths(infos map[string]sharedInfo) (res []string) {
	for _, v := range infos {
		res = append(res, v.path)
	}
	return
}

// LoadTemplates initiates available templates.
func LoadTemplates() {

	if templates == nil {
		loadConfiguration()
		templates = make(map[string]templateInfo)
		layouts = make(map[string]map[string]sharedInfo)
		shared = make([]string, 0, 0)
		sharedCheck = make(map[string]sharedInfo)
	}

	mainTemplate := template.New("main")
	mainTemplate, err := mainTemplate.Parse(mainTmpl)

	err = filepath.Walk(config.MakePath(templateConfig.TemplateLayoutPath), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Templates: prevent panic by handling failure accessing a path %q: %v\n", templateConfig.TemplateLayoutPath, err)
			return err
		}
		if strings.HasSuffix(info.Name(), ".tmpl") {

			layoutfullname := strings.TrimLeft(strings.Replace(path, config.MakePath(templateConfig.TemplateLayoutPath), "", 1), config.SYS_DIR_SEP)
			layoutbase := strings.TrimRight(strings.Replace(layoutfullname, info.Name(), "", 1), config.SYS_DIR_SEP)
			layoutname := strings.TrimLeft(layoutfullname, config.SYS_DIR_SEP)

			var tmpl sharedInfo
			tmpl.path = path
			tmpl.lastUpdate = time.Now().UTC()

			if _, found := layouts[layoutbase][layoutname]; !found {
				layouts[layoutbase] = make(map[string]sharedInfo)
			}

			layouts[layoutbase][layoutname] = tmpl

			log.Printf("Templates: Added Layout of file : %s as %s", path, layoutbase)

		}

		return nil
	})

	if err != nil {
		log.Fatalf("Templates: Failed to load layout templates: %s\n", err)
	}

	err = filepath.Walk(config.MakePath(templateConfig.TemplateSharedPath), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Templates: prevent panic by handling failure accessing a path %q: %v\n", config.MakePath(templateConfig.TemplateSharedPath), err)
			return err
		}

		if strings.HasSuffix(info.Name(), ".tmpl") {
			shared = append(shared, path)
			var sharedinfo sharedInfo
			sharedinfo.path = path
			sharedinfo.lastUpdate = time.Now().UTC()
			sharedCheck[path] = sharedinfo
			log.Printf("Templates: Added shared of file : %s ", path)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Templates: Failed to load shared templates: %s\n", err)
	}

	err = filepath.Walk(config.MakePath(templateConfig.TemplateIncludePath), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Templates: prevent panic by handling failure accessing a path %q: %v\n", config.MakePath(templateConfig.TemplateIncludePath), err)
			return err
		}
		if strings.HasSuffix(info.Name(), ".tmpl") {

			templatefullname := strings.Replace(strings.TrimLeft(strings.Replace(path, config.MakePath(templateConfig.TemplateIncludePath), "", 1), config.SYS_DIR_SEP), ".html.tmpl", "", 1)
			templatebase := strings.Split(templatefullname, config.SYS_DIR_SEP)[0]

			var tmpl templateInfo

			tmpl.baseTmpl, err = mainTemplate.Clone()
			PreLoadFunctions(tmpl.baseTmpl)
			if err != nil {
				log.Fatalf("Templates: Failed to clone mainTemplate: %s\n", err)
			}

			files := append(paths(layouts[""]), append(paths(layouts[templatebase]), path)...)
			tmpl.tmpl = template.Must(tmpl.baseTmpl.ParseFiles(files...))
			tmpl.lastUpdate = time.Now().UTC()
			tmpl.path = path
			tmpl.base = templatebase

			templates[templatefullname] = tmpl

			log.Printf("Templates: Loaded template : %s as %s using layout %s\n", path, templatefullname, templatebase)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Templates: Failed to load templates: %s", err)

	}

	log.Printf("Templates: Loading successful Available: %d: %v", len(templates), reflect.ValueOf(templates).MapKeys())

	bufpool = bpool.NewBufferPool(64)
	log.Println("Templates: buffer allocation successful")
}

// if shared has been updated need to reload all templates ...
func checkShared() {
	tmpShared := make([]string, 0, 0)
	tmpSharedCheck := make(map[string]sharedInfo)
	altered := false
	err := filepath.Walk(config.MakePath(templateConfig.TemplateSharedPath), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Templates: prevent panic by handling failure accessing a path %q: %v\n", config.MakePath(templateConfig.TemplateLayoutPath), err)
			return err
		}

		if strings.HasSuffix(info.Name(), ".tmpl") {
			var shif sharedInfo
			shif.lastUpdate = info.ModTime()
			shif.path = path
			tmpSharedCheck[path] = shif
			tmpShared = append(tmpShared, path)
			sharedinfo, found := sharedCheck[path]
			if !found {
				log.Printf("Templates: Added shared of file : %s ", path)
				altered = true
			} else {
				if info.ModTime().After(sharedinfo.lastUpdate) {
					log.Printf("Templates: Shared file has been altered : %s ", path)
					altered = true
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Templates: Failed to load shared templates: %s\n", err)
	}

	if altered {
		sharedCheck = tmpSharedCheck
		shared = tmpShared
		mainTemplate := template.New("main")
		mainTemplate, _ = mainTemplate.Parse(mainTmpl)

		log.Printf("Templates: Rebuilding templates as shared have evolved...")
		for k, v := range templates {
			files := append(append(paths(layouts[""]), append(paths(layouts[v.base]), v.path)...), shared...)
			v.baseTmpl, err = mainTemplate.Clone()
			PreLoadFunctions(v.baseTmpl)
			v.tmpl = template.Must(v.baseTmpl.ParseFiles(files...))
			v.lastUpdate = time.Now().UTC()
			templates[k] = v
		}
	}

}

// if layouts has been updated need to reload all templates ...
func checkLayouts() {
	tmpLayoutCheck := make(map[string]map[string]sharedInfo)
	altered := false
	err := filepath.Walk(config.MakePath(templateConfig.TemplateLayoutPath), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Templates: prevent panic by handling failure accessing a path %q: %v\n", config.MakePath(templateConfig.TemplateLayoutPath), err)
			return err
		}

		if strings.HasSuffix(info.Name(), ".tmpl") {
			layoutfullname := strings.TrimLeft(strings.Replace(path, config.MakePath(templateConfig.TemplateLayoutPath), "", 1), config.SYS_DIR_SEP)
			layoutbase := strings.TrimRight(strings.Replace(layoutfullname, info.Name(), "", 1), config.SYS_DIR_SEP)
			layoutname := strings.TrimLeft(layoutfullname, config.SYS_DIR_SEP)

			var shif sharedInfo
			shif.lastUpdate = info.ModTime()
			shif.path = path

			if _, found := tmpLayoutCheck[layoutbase][layoutname]; !found {
				tmpLayoutCheck[layoutbase] = make(map[string]sharedInfo)
			}
			tmpLayoutCheck[layoutbase][layoutname] = shif
			// iterate on all layouts ...
			_, found := layouts[layoutbase]
			if !found {
				log.Printf("Templates: Added Layout of file : %s ", layoutfullname)
				altered = true
			} else {
				locallayout, found := layouts[layoutbase][layoutname]
				if !found {
					log.Printf("Templates: Added Layout of file : %s ", layoutfullname)
					altered = true
				} else {
					if shif.lastUpdate.After(locallayout.lastUpdate) {
						log.Printf("Templates: Layout file has been altered : %s ", layoutfullname)
						altered = true
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Templates: Failed to load shared templates: %s\n", err)
	}

	if altered {
		layouts = tmpLayoutCheck
		mainTemplate := template.New("main")
		mainTemplate, _ = mainTemplate.Parse(mainTmpl)

		log.Printf("Templates: Rebuilding templates as shared have evolved...")
		for k, v := range templates {
			files := append(append(paths(layouts[""]), append(paths(layouts[v.base]), v.path)...), shared...)
			v.baseTmpl, err = mainTemplate.Clone()

			PreLoadFunctions(v.baseTmpl)
			v.tmpl = template.Must(v.baseTmpl.ParseFiles(files...))
			v.lastUpdate = time.Now().UTC()
			templates[k] = v
		}
	}

}

//RenderTemplateFn with custom functions provided.
// Be ware these functions are only valid with this call !
func RenderTemplateFn(w http.ResponseWriter, req *http.Request, name string, data interface{}, fns template.FuncMap) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		log.Fatalf("Templates: The template %s does not exist. Can't render. Available: %d: %v", name, len(templates), reflect.ValueOf(templates).MapKeys())
		return
	}

	if config.WEB_RELOADING {
		// reload shared stuff.
		checkLayouts()
		checkShared()

		// seek last update ...
		file, err := os.Open(tmpl.path)
		if err != nil {
			http.Error(w, "Failed to render page - page missing", http.StatusInternalServerError)
			log.Fatalf("Templates: The template %s does not exist. Can't render.", name)
			return
		}

		info, _ := file.Stat()

		if info.ModTime().After(tmpl.lastUpdate) {
			log.Printf("Templates: An update is available for template: %s - %s", name, tmpl.path)
			mainTemplate := template.New("main")
			mainTemplate, _ = mainTemplate.Parse(mainTmpl)
			tmpl.baseTmpl, err = mainTemplate.Clone()
			PreLoadFunctions(tmpl.baseTmpl)
			files := append(append(paths(layouts[""]), append(paths(layouts[tmpl.base]), tmpl.path)...), shared...)
			tmpl.tmpl = template.Must(tmpl.baseTmpl.ParseFiles(files...))
			tmpl.lastUpdate = time.Now().UTC()
			templates[name] = tmpl
		}
	}

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	LoadFunctions(w, req, tmpl.tmpl, fns)

	err := tmpl.tmpl.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := tools.GetSession(req).Save(req, w); err != nil {
		log.Fatalf("Error saving session: %v", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)

}

//RenderTemplate render provided templates name. Template name must match path eg: garden/index
func RenderTemplate(w http.ResponseWriter, req *http.Request, name string, data interface{}) {
	fn := make(template.FuncMap)

	RenderTemplateFn(w, req, name, data, fn)
}
