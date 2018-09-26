package networks

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type WebAcceptLanguage struct {
	Language string
	Quality  float32
}

type WebAcceptLanguages []WebAcceptLanguage

func (l WebAcceptLanguages) Len() int           { return len(l) }
func (l WebAcceptLanguages) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l WebAcceptLanguages) Less(i, j int) bool { return l[i].Quality > l[j].Quality }
func (l WebAcceptLanguages) String() string {
	output := bytes.NewBufferString("")
	for i, language := range l {
		output.WriteString(fmt.Sprintf("%s (%1.1f)", language.Language, language.Quality))
		if i != len(l)-1 {
			output.WriteString(", ")
		}
	}

	return output.String()
}

func ResolveAcceptLanguage(req *http.Request) WebAcceptLanguages {
	header := req.Header.Get("Accept-Language")
	if header == "" {
		return nil
	}

	acceptLanguageHeaderValues := strings.Split(header, ",")
	acceptLanguages := make(WebAcceptLanguages, len(acceptLanguageHeaderValues))

	for i, languageRange := range acceptLanguageHeaderValues {
		if qualifiedRange := strings.Split(languageRange, ";q="); len(qualifiedRange) == 2 {
			quality, err := strconv.ParseFloat(qualifiedRange[1], 32)
			if err != nil {
				acceptLanguages[i] = WebAcceptLanguage{qualifiedRange[0], 1}
			} else {
				acceptLanguages[i] = WebAcceptLanguage{qualifiedRange[0], float32(quality)}
			}
		} else {
			acceptLanguages[i] = WebAcceptLanguage{languageRange, 1}
		}
	}

	sort.Sort(acceptLanguages)
	return acceptLanguages
}
