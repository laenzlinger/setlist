package template

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/laenzlinger/setlist/internal/config"
)

var (
	//go:embed *.html
	templateFS embed.FS

	//nolint:gochecknoglobals // we want to check the templates on application start
	placeholderTemplate *template.Template
	//nolint:gochecknoglobals // we want to check the templates on application start
	setlistTemplate *template.Template
)

//nolint:gochecknoinits // we want to check the templates on application start
func init() {
	setlistTemplate = template.Must(template.New("setlist.html").ParseFS(templateFS, "setlist.html"))
	placeholderTemplate = template.Must(template.New("placeholder.html").ParseFS(templateFS, "placeholder.html"))
}

type Data struct {
	Title   string
	Margin  string
	Content template.HTML
}

func CreateSetlist(data *Data) (string, error) {
	filename := filepath.Join(config.Target(), fmt.Sprintf("Set List %s.html", data.Title))
	return createFromTemplate(data, setlistTemplate, filename)
}

func CreatePlaceholder(data *Data) (string, error) {
	return createFromTemplate(data, placeholderTemplate, filepath.Join(config.PlaceholderDir(), "placeholder.html"))
}

func createFromTemplate(data *Data, t *template.Template, filename string) (string, error) {
	PrepareTarget()
	f, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create HTML file: %w", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Println(err)
		}
	}()

	err = t.Execute(f, data)
	if err != nil {
		return "", fmt.Errorf("failed to exectue HTML template: %w", err)
	}
	return filename, nil
}

func PrepareTarget() {
	if err := os.MkdirAll(config.Target(), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
