package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Структуры GitHub API
type IssuesSearchResult struct {
	TotalCount int     `json:"total_count"`
	Items      []*Issue `json:"items"`
}

type Issue struct {
	Number    int       `json:"number"`
	HTMLURL   string    `json:"html_url"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	User      *User     `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	Body      string    `json:"body"`
}

type User struct {
	Login   string `json:"login"`
	HTMLURL string `json:"html_url"`
}

// Шаблон
const templ = `{{.TotalCount}} тем найдено в GitHub:

{{range .Items}}----------------------------------------
📌 #{{.Number}} [{{.State | upper}}] {{.Title | printf "%.65s"}}
   👤 {{.User.Login}} | 📅 {{.CreatedAt | daysAgo}} дней назад
   🔗 {{.HTMLURL}}

{{end}}`

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

func upper(s string) string {
	return strings.ToUpper(s)
}

var report = template.Must(template.New("issuelist").
	Funcs(template.FuncMap{
		"daysAgo": daysAgo,
		"printf":  fmt.Sprintf,
		"upper":   upper,
	}).
	Parse(templ))

const IssuesURL = "https://api.github.com/search/issues"

func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	
	client := &http.Client{Timeout: 30 * time.Second}
	
	resp, err := client.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сбой запроса: %s", resp.Status)
	}
	
	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Использование: %s <поисковые_термы>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Примеры:\n")
		fmt.Fprintf(os.Stderr, "  %s repo:golang/go is:open\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s language:go stars:>1000\n", os.Args[0])
		os.Exit(1)
	}
	
	result, err := SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	
	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}