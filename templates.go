package main

import "html/template"

func loadTemplates() (*template.Template, error) {
	t := template.New("")
	if _, err := t.New("login").Parse(loginTmpl); err != nil {
		return nil, err
	}
	if _, err := t.New("dash").Parse(dashTmpl); err != nil {
		return nil, err
	}
	return t, nil
}
