## Article and Taggin, built in Golang
---

#### Requirements
- github: Installtion [instructions](https://github.com/git-guides/install-git)
- docker: Download from [here](https://docs.docker.com/get-docker/)

---
### How to run
1. Git clone this repo
2. Open terminal/shell and change directory to cloned directory
3. Make sure docker is running
4. Run `docker-compose up --build`

The docker will download images and build containers. Once it has completed, you can access the APIs on http://localhost:8080/

---
### Questions
1. Port 8080 taken?

   Change the "APP_PORT" from .env file in project home directory

2. Permission error on docker entrypoint bash scripts?

   As docker inherits permissions from your system, you need to change their permission to executable. Run `sudo chmod 775 script.sh` on the file
---
### Summary

The project is built in Golang. I am a Python/Django developer, however, I went with Golang because -

1. It's fast. We all know it
2. I like the language and have been meaning to look for an excuse to work in it
3. Extra bonus points

I have tried following Golang project best practices mentioned [here](https://peter.bourgon.org/go-best-practices-2016/?ref=hackernoon.com#repository-structure). 

---
### Assumptions
1. The date in article field is assumed to be post creation date. Therefore, it is generated on backend and is not accepted as a json field
2. Duplicate article title and article body are not allowed
---
### How to add a new article
Make a post request on http://localhost:8080/articles with application/json in request body that follows the following structure

```
{
    "articleTitle":"Golang is Awesome",
    "articleBody":"article body text",
    "tags":["tag3", "tag4", "tag1", "tag2"]
}
```
---
### If I had more time
These are the following TODOs that I had when I finished the project -
1. Add a project-wide logger
2. There is some code repetition. Remove that
3. **Write tests**. Unfortunatly, due to the time crunch, I could not write unit tests for the project. 

---
### How long did it take
It took me approximately 2 days to finish this project. I have worked with ORMs and they do all the heavy lifting for you. I also had limited experience with Go, therefore, I learnt quite a lot on the fly. 

I had the option to use ORMs but I did not as it is anti-pattern in Go and I wanted to keep my learning experience as much as possible. Oh, I also implemented the same project in Python/Django to compare my learning curve and time taken. I completed the project in Python/Django in 5 hours. However, I am presenting Go version of the project to you. This was challenging and rewarding.

>  Thanks for going through the project :)


