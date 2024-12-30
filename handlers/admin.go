package handlers

import (
	"html/template"
	"net/http"

	"github.com/pranavKharche24/mail/utils"
)

func HandleAdmin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	data := struct {
		FromEmail string
		FromPass  string
	}{
		FromEmail: utils.FromEmail,
		FromPass:  utils.FromPass,
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func HandleAdminSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	utils.FromEmail = r.FormValue("fromEmail")
	utils.FromPass = r.FormValue("fromPass")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
