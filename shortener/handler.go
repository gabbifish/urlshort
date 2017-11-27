package shortener

import (
	"net/http"
)

func handleHandlerFuncInternalError(w http.ResponseWriter, err error, errMsg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(errMsg))
}

func NewHandlerFromSlugToURLClient(mapper SlugToURL, fallback http.Handler) http.Handler {
	// TODO: What do we do when there isn't a fallback handler?
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Does a URL exist for the given slug?
		slug := r.URL.Path
		slugExists, slugExistsErr := mapper.URLExists(slug)
		if slugExistsErr != nil {
			handleHandlerFuncInternalError(w, slugExistsErr, "Oops!  Our slug-to-url client threw an error while checking for a matching URL")
			return
		}
		// If a URL does not exist, then do whatever the fallback handler says
		if !slugExists {
			fallback.ServeHTTP(w, r)
			return
		}
		// Get the appropriate URL, given the slug
		url, getURLFromSlugErr := mapper.URL(r.URL.Path)
		if getURLFromSlugErr != nil {
			handleHandlerFuncInternalError(w, getURLFromSlugErr, "Oops!  Our slug-to-url client threw an error while getting the matching URL")
			return
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	})
}
