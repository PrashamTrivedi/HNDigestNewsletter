package hndata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

//FetchHackernewsData fetches hackernews data from given type
func FetchHackernewsData(argsType string) (string, string) {

	var story, storyType string
	switch argsType {
	case "ask":
		story = "askStories"
		storyType = "Ask HN"
	case "show":
		story = "showstories"
		storyType = "Show HN"
	case "job":
		story = "jobstories"
		storyType = "Jobs"
	default:
		story = "topstories"
		storyType = "Top Stories"
	}
	fmt.Println("Story is", story, "Story Type is", storyType)
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/%v.json", story)

	var storyIds []int
	storiesData := getData(url, &storyIds)
	if storiesData.errorData != nil {
		panic(storiesData.errorData)
	}

	fmt.Println(storyIds)

	stories := make([]Story, 0)
	storiesChannel := make(chan Story, 10)

	for _, storyID := range storyIds[:10] {

		var story Story
		storyURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json?print=pretty", storyID)
		storyData := getData(storyURL, &story)
		if storyData.errorData != nil {
			fmt.Println(storyData.errorData)
		}
		story.TimeString = time.Unix(story.Time, 0).Format("02 Nov 2006")
		if story.Kids != nil {
			for _, kidID := range story.Kids[:3] {
				if kidID != 0 {

					var kid Story
					kidURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json?print=pretty", kidID)
					kidData := getData(kidURL, &kid)
					if kidData.errorData != nil {
						fmt.Println(kidData.errorData)
					}
					kid.Text = strings.SplitN(kid.Text, "<p>", 2)[0]
					kid.TimeString = time.Unix(kid.Time, 0).Format("02 Nov 2006")
					story.KidData = append(story.KidData, kid)
				}
			}
		}
		// stories[index] = story
		storiesChannel <- story

	}
	close(storiesChannel)

	for data := range storiesChannel {
		// fmt.Println(data)
		stories = append(stories, data)
	}

	htmlTemplate, err := template.New("HTML Template").Funcs(template.FuncMap{
		"htmlSafe": func(html string) template.HTML {
			return template.HTML(html)
		},
	}).Parse(`<html>
            <head>
                <meta charset="utf-8">
                <meta name="HandheldFriendly" content="True">
                <meta name="MobileOptimized" content="320">
                <meta name="viewport" content="width=device-width, initial-scale=1.0">
                <meta name="referrer" content="no-referrer">
                <style>
                    * {
                        border:0;
                        font:inherit;
                        font-size:100%;
                        vertical-align:baseline;
                        margin:0;
                        padding:0;
                        color: black;
                        text-decoration-skip: ink;
                    }

                    h1 {
                        font-size:25px;
                        margin-top:18px;
                    }

                    h2 {
                        font-size:18px;
                    }

                    h3 {
                        font-size:12px;
                    }

                    h1 a,h2 a,h3 a {
                        text-decoration:none;
                    }
                    body {
                        font-family:'Roboto' , sans-serif;
                        font-size:14px;
                        color:#1d1313;
                        max-width:700px;
                        margin:auto;
                    }

                    .commentBox {
                        padding: 0px 8px;
                        margin-bottom: 5px;
                        background:#403c3b;
                    }

                    @media (prefers-color-scheme: dark) {
                            *, #nav h1 a {
                                color: #FDFDFD;
                            }

                            body {
                                background: #121212;
                            }

                            pre, code {
                                background-color: #262626;
                            }

                            #sub-header, .date {
                                color: #BABABA;
                            }

                            hr {
                                background: #EBEBEB;
                            }
                        }
                </style>
            </head>
            <body>
				Here is your Digest
				{{range .}}
					<h1>
						{{if .URL}}
							<a href="{{.URL}}">{{.Title}}</a>
						{{else}}
							{{.Title}}
						{{end}}
					</h1>
					<p> At {{.TimeString}}
						<p> By 
						<a href="https://news.ycombinator.com/user?id=${{.By}}" target="_blank">{{.By}}</a>
					<p> {{.Text}}

					{{if .KidData}}
						<h2>Comments</h2>
						{{ range .KidData }}
						<div class="commentBox">

							{{if .URL}}
								<a href="{{.URL}}">{{.Title}}</a>
							{{else}}
								{{.Title}}
							{{end}}
							<p> At {{.TimeString}}
							<p> By 
							<a href="https://news.ycombinator.com/user?id={{.By}}" target="_blank">{{.By}}</a>
							<p> {{.Text | htmlSafe}}
							   <a href="https://news.ycombinator.com/item?id={{.ID}}" target="_blank">Read More...</a>
						</div>
						<p> 
						{{end}}
					{{end}}
				{{end}}
                
            </body>
			</html>`)

	if err != nil {
		panic(err)
	}
	htmlFile, fileError := os.Create("../testData/fbData_go.html")

	if fileError != nil {
		panic(fileError)
	}

	returnData := ""
	env := os.Getenv("ENV")
	if env != "" && env == "PROD" {
		var templateData bytes.Buffer
		if err := htmlTemplate.Execute(&templateData, stories); err != nil {
			panic(err)
		}

		returnData = templateData.String()
	} else {

		templateError := htmlTemplate.Execute(htmlFile, stories)
		if templateError != nil {
			panic(templateError)
		}
	}
	return returnData, storyType
}

func getData(url string, typeData interface{}) APIResponse {
	fmt.Println(url)
	apiResponse := APIResponse{}
	response, error := http.Get(url)
	if error != nil {
		apiResponse.errorData = error
	} else {
		defer response.Body.Close()
		err := json.NewDecoder(response.Body).Decode(&typeData)
		ioutil.ReadAll(response.Body)
		if err != nil {
			apiResponse.errorData = err
		} else {

			apiResponse.data = typeData
			apiResponse.errorData = nil

		}
	}
	return apiResponse

}
