package funcs

import "net/http"

func AboutHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	RenderTemplate(w, "about.html", nil)

}
