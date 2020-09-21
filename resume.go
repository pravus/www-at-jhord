package main

import (
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "text/template"

  "github.com/go-chi/chi"
  "gopkg.in/yaml.v2"
)

type Resume struct {
  Contact struct {
    FirstName string   `json:"first_name" yaml:"first_name"`
    LastName  string   `json:"last_name"  yaml:"last_name"`
    Title     string   `json:"title"      yaml:"title"`
    Email     string   `json:"email"      yaml:"email"`
    Phone     string   `json:"phone"      yaml:"phone"`
    Links     []struct {
      Name string `json:"name" yaml:"name"`
      Href string `json:"href" yaml:"href"`
    } `json:"links" yaml:"links"`
    Address struct {
      Line1     string `json:"line1"   yaml:"line1"`
      Line2     string `json:"line2"   yaml:"line2"`
      City      string `json:"city"    yaml:"city"`
      State     string `json:"state"   yaml:"state"`
      ZipCode   string `json:"zipcode" yaml:"zipcode"`
    }                  `json:"address" yaml:"address"`
  } `json:"contact" yaml:"contact"`
  Employments []struct {
    Name       string   `json:"name"       yaml:"name"`
    Title      string   `json:"title"      yaml:"title"`
    Years      string   `json:"years"      yaml:"years"`
    City       string   `json:"city"       yaml:"city"`
    State      string   `json:"state"      yaml:"state"`
    Keywords   []string `json:"keywords"   yaml:"keywords"`
    Experience []string `json:"experience" yaml:"experience"`

  } `json:"employments" yaml:"employments"`
  Skills []string `json:"skills" yaml:"skills"`
  Education struct {
    Institution string `json:"institution" yaml:"institution"`
    Degree      string `json:"degree"      yaml:"degree"`
    Major       string `json:"major"       yaml:"major"`
    Minor       string `json:"minor"       yaml:"minor"`
    Years       string `json:"years"       yaml:"years"`
    Status      string `json:"status"      yaml:"status"`
  } `json:"education" yaml:"education"`
}

type Style struct {
  Palette struct {
    Background string `json:"background" yaml:"background"`
    Text       string `json:"text"       yaml:"text"`
    Borders    string `json:"borders"    yaml:"borders"`
    Hilite     string `json:"hilite"     yaml:"hilite"`
  } `json:"palette" yaml:"palette"`
}

func ResumeRouter() func (chi.Router) {
  return func (r chi.Router) {
    r.Get("/",    ResumeHTML())
    r.Get("/css", ResumeCSS ())
  }
}

func ResumeHTML() http.HandlerFunc {
  resume  := ResumeLoadFromYaml("resume/content.yaml")
  content := template.Must(template.ParseFiles("resume/content.html"))
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/html")
    err := content.Execute(w, resume)
    if err != nil {
      log.Printf("template: error=%v", err)
      return
    }
  })
}

func ResumeCSS() http.HandlerFunc {
  content := template.Must(template.ParseFiles("resume/content.css"))
  style   := StyleLoadFromYaml("resume/style.yaml")
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/css")
    err := content.Execute(w, style)
    if err != nil {
      log.Printf("template: error=%v", err)
    }
  })
}

func ResumeLoadFromYaml(name string) *Resume {
  file, err := os.Open(name)
  if err != nil {
    log.Panicf("open: %s: %v", name, err)
  }
  defer file.Close()

  body, err := ioutil.ReadAll(file)
  if err != nil {
    log.Panicf("read: %s: %v", name, err)
  }

  resume := &Resume{}
  err = yaml.Unmarshal(body, resume)
  if err != nil {
    log.Panicf("yaml: %s: %v", name, err)
  }

  return resume
}

func StyleLoadFromYaml(name string) *Style {
  file, err := os.Open(name)
  if err != nil {
    log.Panicf("open: %s: %v", name, err)
  }
  defer file.Close()

  body, err := ioutil.ReadAll(file)
  if err != nil {
    log.Panicf("read: %s: %v", name, err)
  }

  style := &Style{}
  err = yaml.Unmarshal(body, style)
  if err != nil {
    log.Panicf("yaml: %s: %v", name, err)
  }

  return style
}
