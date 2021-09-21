package server

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	. "github.com/guricerin/stop-now-smoking/util"
)

var (
	funcMap = map[string]interface{}{
		"isLoginAndRsrcUserMatch": func(state LoginState) bool {
			return state == LoginAndRsrcUser
		},
		"now": func() string {
			return time.Now().Format(timeLayout)
		},
		"dateFormat": func(t time.Time) string {
			return t.Format(timeLayout)
		},
	}
)

func writeHtml(w http.ResponseWriter, viewModel interface{}, filenames ...string) (err error) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.New("funcMap").Funcs(funcMap).ParseFiles(files...))
	err = templates.ExecuteTemplate(w, "layout", viewModel)
	if err != nil {
		Elog.Printf("execute template failed: %v", err)
		fmt.Fprintf(w, "%v", err)
	}
	return
}
