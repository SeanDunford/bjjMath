package UrlResolver

import (
	"log"
	"net/http"
)

func ResolveUrl(unresolvedUrl string) string {
	resp, err := http.Get(unresolvedUrl)
	if err != nil {
		log.Fatalf("http.Get => %v", err.Error())
	}

	resolvedUrl := resp.Request.URL.String()
	return resolvedUrl
}
