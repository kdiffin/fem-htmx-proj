package main

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"

	"myapp/database"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("../views/*.gohtml")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Renderer = newTemplate()
	e.Use(middleware.Logger())
	e.Static("/images", "../images")
	e.Static("/css", "../css")
	e.Static("/js", "../js")

	dbpool := database.DatabaseInit()
	// close db pool when the funciton main ends
	defer dbpool.Close()

	// todo: reimplement
	// router.Tasks(e, dbpool)
	// router.Tags(e, dbpool)
	e.GET("/", func(c echo.Context) error {

		rows, err := dbpool.Query(context.Background(), "SELECT (name, idea, id, tags, completed) FROM tasks")
		if err != nil {
			c.String(500, http.StatusText(500))
		}

		type IndexPage struct {
			Tasks database.Tasks
		}
		indexPage := IndexPage{
			Tasks: database.Tasks{},
		}

		for rows.Next() {
			var Task database.Task
			errScan := rows.Scan(&Task)
			if errScan != nil {
				c.String(500, http.StatusText(500))
			}

			indexPage.Tasks = append(indexPage.Tasks, Task)
		}

		return c.Render(200, "index", indexPage)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
