package template

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"os"
)

var (
	//go:embed *.html
	templateFS embed.FS

	placeholderTemplate *template.Template
	setlistTemplate     *template.Template
)

func init() {
	setlistTemplate = template.Must(template.New("setlist.html").ParseFS(templateFS, "setlist.html"))
	placeholderTemplate = template.Must(template.New("placeholder.html").ParseFS(templateFS, "placeholder.html"))
}

type TemplateData struct {
	Title   string
	Content template.HTML
}

func CreateSetlist(data *TemplateData) string {
	return createFromTemplate(data, setlistTemplate, fmt.Sprintf("out/Setlist %s.html", data.Title))
}

func CreatePlaceholder(data *TemplateData) string {
	return createFromTemplate(data, placeholderTemplate, "out/placeholder/placeholder.html")
}

func createFromTemplate(data *TemplateData, t *template.Template, filename string) string {
	PrepareOut()
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	err = t.Execute(f, data)
	if err != nil {
		log.Fatal(err)
	}
	return filename
}

func PrepareOut() {
	if err := os.MkdirAll("out", os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
