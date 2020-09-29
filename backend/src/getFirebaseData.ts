import https from 'https'
import fs from 'fs'


export const getFirebaseData = async (event: any) => {
    try {
        console.log(event)
        let story = ''
        if (!event) {
            story = 'topstories'
        } else if (event.type === 'ask') {
            story = 'askstories'
        } else if (event.type === 'show') {
            story = 'showstories'
        } else if (event.type === 'job') {
            story = 'jobstories'
        } else {
            story = 'topstories'
        }
        const askStoriesData = await getData(`https://hacker-news.firebaseio.com/v0/${story}.json`)
        // console.log(askStoriesData)
        const stories = []

        for await (const storyId of askStoriesData.slice(0, 10)) {

            const storyData = await getData(`https://hacker-news.firebaseio.com/v0/item/${storyId}.json?print=pretty`)
            const date = new Date()
            date.setSeconds(storyData.time)
            storyData.time = date.toISOString().replace('T', ' ').replace('Z', '')
            if (storyData.kids) {
                const kids = []

                for await (const kidId of storyData.kids.slice(0, 3)) {

                    const kidData = getData(`https://hacker-news.firebaseio.com/v0/item/${kidId}.json?print=pretty`).then(kid => {
                        if (kid && !kid.deleted && kid.time) {
                            const kidDate = new Date()
                            kidDate.setSeconds(kid.time)
                            kid.time = kidDate.toISOString().replace('T', ' ').replace('Z', '')
                            return kid
                        }
                    }).catch(error => console.log(error))
                    kids.push(kidData)


                }
                storyData.kidData = await Promise.all(kids)
            }

            stories.push(storyData)
        }

        const storiesInHtml = stories.map(story => {
            const text = story.text || ''

            return `
            <h1>
                ${story.url ?
                    `<a href="${story.url}" target="_blank">${story.title}</a>` : `${story.title}`}
            </h1>
            <p> At ${story.time}
            <p> By <a href="https://news.ycombinator.com/user?id=${story.by}" target="_blank">${story.by}</a>
            <p> ${text}
            
            ${story.kidData ? `<h2>Comments</h2>` : ''}
            ${story.kidData && story.kidData.filter((kid: any) => kid).map((kid: {time: any; by: string; text: string, id: number}) => {
                        const length = kid.text.indexOf('<p>') === -1 ? 20 : Math.max(50, kid.text.indexOf('<p>'))
                        return `  
                <div class="commentBox">
                <p> At ${kid.time}
                <p> By <a href="https://news.ycombinator.com/user?id=${kid.by}" target="_blank">${kid.by}</a>
                <p> ${kid.text && kid.text.substring(0, length)}
                <a href="https://news.ycombinator.com/item?id=${kid.id}" target="_blank">Read More...</a>
                </div>
                `
                    }).join(`<p>      `)}
            ${story.kidData ? `<a href="https://news.ycombinator.com/item?id=${story.id}" target="_blank">All Comments</a>` : ''}
           
            `}

        )
        const htmlData = `
         <html>
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
                ${storiesInHtml.join('<p>').split('undefined').join(' ')}
            </body>
            </html>`

        if (!fs.existsSync('testData')) {
            fs.mkdirSync('testData')
        }
        fs.writeFileSync('./testData/fbData.html', htmlData)
        console.log(JSON.stringify(stories))
    } catch (error) {
        console.log(error)
    }
}

async function getData(url: string | https.RequestOptions | URL): Promise<any> {
    return new Promise((resolve, reject) => {
        https.get(url, (res) => {

            res.on('data', (d) => {
                resolve(JSON.parse(d.toString()))
            })

        }).on('error', (e) => {
            reject(e)
        })
    })

}

