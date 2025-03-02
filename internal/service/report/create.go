package report

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"diploma-project/internal/model"

	"github.com/google/uuid"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

var HEADERS = map[string]string{
	"accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	"user-agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
}

func (s *Service) Create(ctx context.Context, cmd model.ReportCreateCommand) error {
	client := &http.Client{}

	u, err := url.Parse(cmd.URL)
	if err != nil {
		return fmt.Errorf("parse url: %w: %w", err, model.ErrNotValid)
	}

	host := u.Scheme + "://" + u.Host

	req, _ := http.NewRequest("GET", cmd.URL, nil)
	for k, v := range HEADERS {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w: %w", err, model.ErrNotValid)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("parse response body: %w", err)
	}

	cmd.Report = htmlCheck(doc, cmd.URL)
	scriptsCheck(getScripts(doc, host))

	errg, gctx := errgroup.WithContext(ctx)

	errg.Go(func() error {
		return s.r.Create(gctx, cmd)
	})

	errg.Go(func() error {
		return s.s.Create(gctx, model.SQLMapCommand{
			ID:     cmd.Report.Id,
			URL:    cmd.URL,
			Report: nil,
		})
	})

	err = errg.Wait()
	if err != nil {
		return fmt.Errorf("reports saving: %w", err)
	}

	return nil
}

func getScripts(doc *html.Node, host string) []string {
	var scripts []string
	var srcList []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					srcList = append(srcList, attr.Val)
				}
			}
			scripts = append(scripts, renderNode(n))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	file, _ := os.Create("scripts.js")
	defer file.Close()

	for _, script := range scripts {
		file.WriteString(script + "\n\n\n")
	}

	for _, src := range srcList {
		resp, err := http.Get(host + src)
		if err == nil {
			body, _ := ioutil.ReadAll(resp.Body)
			file.WriteString(string(body) + "\n\n\n")
			scripts = append(scripts, string(body))
			resp.Body.Close()
		}
	}

	return scripts
}

func htmlCheck(doc *html.Node, source string) model.Report {
	var report model.Report
	report.Id = uuid.New().String()

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "form" {
				method := getAttr(n, "method")
				uri := getAttr(n, "action")
				if uri == "" {
					uri = "/"
				}

				u, _ := url.Parse(uri)
				url2, _ := url.Parse(source)
				if u.Host != "" && u.Host != url2.Host {
					return
				}

				keys := make(map[string][]string)
				var g func(*html.Node)
				g = func(n *html.Node) {
					if n.Type == html.ElementNode && n.Data == "input" {
						name := getAttr(n, "name")
						if name != "" {
							keys[name] = []string{}
							for comment, r := range inputRegexps {
								if match, _ := regexp.MatchString(r, renderNode(n)); match {
									keys[name] = append(keys[name], comment)
								}
							}
						}
					}
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						g(c)
					}
				}
				g(n)

				params := make([]model.Param, 0, len(keys))
				for name, comments := range keys {
					params = append(params, model.Param{
						Name:     name,
						Patterns: comments,
					})
				}

				report.URLs = append(report.URLs, model.URL{
					URL:    uri,
					Method: method,
					Params: params,
				})
			} else if n.Data == "a" {
				uri := getAttr(n, "href")
				if uri == "" {
					return
				}

				u, _ := url.Parse(uri)
				url2, _ := url.Parse(source)
				if u.Host != "" && u.Host != url2.Host {
					return
				}

				params := []model.Param{}
				queryParams := u.Query()
				for name, values := range queryParams {
					var patterns []string
					for _, value := range values {
						for comment, r := range aRegexps {
							if match, _ := regexp.MatchString(r, value); match {
								patterns = append(patterns, comment)
							}
						}
					}

					params = append(params, model.Param{
						Name:     name,
						Values:   values,
						Patterns: patterns,
					})
				}

				report.URLs = append(report.URLs, model.URL{
					URL:    uri,
					Method: "GET",
					Params: params,
				})
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return report
}

func scriptsCheck(scripts []string) {
	file, _ := os.Create("log.txt")
	defer file.Close()

	for _, script := range scripts {
		for _, r := range scriptRegexps {
			if match, _ := regexp.MatchString(r, script); match {
				file.WriteString(script + "\n")
			}
		}
	}
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func renderNode(n *html.Node) string {
	var b strings.Builder
	html.Render(&b, n)
	return b.String()
}
