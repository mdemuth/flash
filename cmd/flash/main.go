package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/mdemuth/flash"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /", http.HandlerFunc(get))
	mux.Handle("POST /", http.HandlerFunc(post))
	if err := http.ListenAndServe(":3333", mux); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	msgs := flash.Get(w, r)

	tmpl, err := template.New("page").Parse(page)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, map[string]any{
		"Flash": msgs,
	})
	if err != nil {
		panic(err)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	s := flash.Severity(r.FormValue("severity"))
	t := r.FormValue("title")
	b := r.FormValue("body")

	flash.Set(w, flash.Message{
		Severity: s,
		Title:    t,
		Body:     b,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

const page = `
<!DOCTYPE html>
<html lang="en">
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>Flash Messages</title>
	<style>
		*, *::before, *::after {
			box-sizing: border-box;
		}

		form {
			border: 1px solid;
			padding: 1em;
		}

		ul {
			list-style-type: none;
			padding-left: 0;
		}

		.flash.message{
			border: 1px solid;
			color: white;
			padding: 1em;
			animation: fade 1s linear 3s;
			animation-fill-mode: forwards; 
		}

		@keyframes fade {
			to { 
				opacity: 0;
			}
		}

		.notice{
			background-color: #333333;
		}

		.info{
			background-color: #6DAAE0;
		}

		.ok{
			background-color: #79A548;
		}

		.warning{
			background-color: #E8A33D;
		}

		.error{
			background-color: #C83C3C;
		}
	</style>
	</head>
	<body>
		<main>
			<h1>Flash Messages</h1>
			<a href="/">Reload</a>
			<form method="POST">
				<p>Set new flash message</p>
				<ul>
					<li>
						<label for="severity">Severity:</label>
						<select name="severity" id="severity">
							<option value="notice">Notice</option>
							<option value="info">Info</option>
							<option value="ok">OK</option>
							<option value="warning">Warning</option>
							<option value="error">Error</option>
						</select>
					</li>
					<li>
						<label for="title">Title:</label>
						<input type="text" id="title" name="title" />
					</li>
					<li>
						<label for="body">Body:</label>
						<input type="text" id="body" name="body" />
					</li>
				</ul>
				<input type="submit" value="set">
			</form>
			{{if .Flash}}
			<ul class="flash messages">
				{{range .Flash}}
					<li class="flash message {{.Severity}}">
						<h3>{{.Title}}</h3>
						<p>{{.Body}}</p>
					</li>
				{{end}}
			</ul>
			{{end}}
		</main>
	</body>
</html>
`
