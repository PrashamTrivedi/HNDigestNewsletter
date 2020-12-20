package hndata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
	// "os"
	// "strings"
	// "time"
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

	stories := make([]Story, 0)

	storyChannels := make(chan int)
	storiesChannel := make(chan Story)

	go getStoryIds(story, storyChannels)

	for storyID := range storyChannels {
		go getStory(storyID, storiesChannel)
	}

	for story := range storiesChannel {
		stories = append(stories, story)
		if len(stories) >= 10 {
			close(storiesChannel)
		}
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

	returnData := ""
	env := os.Getenv("ENV")
	if env == "PROD" {
		var templateData bytes.Buffer
		if err := htmlTemplate.Execute(&templateData, stories); err != nil {
			panic(err)
		}

		returnData = templateData.String()
	} else {
		htmlFile, fileError := os.Create("fbData_go.html")

		if fileError != nil {
			panic(fileError)
		}
		templateError := htmlTemplate.Execute(htmlFile, stories)
		if templateError != nil {
			panic(templateError)
		}
	}
	return returnData, storyType
}

func getStory(storyID int, storyChannel chan Story) {
	var story Story
	storyURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json?print=pretty", storyID)
	if err := getData(storyURL, &story); err != nil {
		panic(err)
	}
	story.TimeString = time.Unix(story.Time, 0).Format("02 Nov 2006")
	populateKids(story, storyChannel)
}

func populateKids(story Story, storyChannel chan Story) {
	if len(story.Kids) > 0 {
		for _, kidID := range story.Kids[:3] {
			if kidID != 0 {
				fmt.Printf("Getting Story for kid %d\n", kidID)
				var kid Story
				kidURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json?print=pretty", kidID)

				if err := getData(kidURL, &kid); err != nil {
					panic(err)
				}
				kid.Text = strings.SplitN(kid.Text, "<p>", 2)[0]
				kid.TimeString = time.Unix(kid.Time, 0).Format("02 Nov 2006")

				story.KidData = append(story.KidData, kid)
			}
		}

		storyChannel <- story
	} else {
		storyChannel <- story
	}
	//close(storyChannel)
}

func getStoryIds(story string, storiesChannel chan int) {
	var storyIds []int
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/%v.json", story)
	if err := getData(url, &storyIds); err != nil {
		panic(err)
	}
	for _, storyId := range storyIds[:10] {
		storiesChannel <- storyId
	}
	close(storiesChannel)
}

func getData(url string, typeData interface{}) error {
	fmt.Println(url)

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&typeData); err != nil {
		return err
	}

	return nil

}
