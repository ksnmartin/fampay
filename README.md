<h1>Backend Assignment | Fampay </h1>

To make an API to fetch latest videos sorted in reverse chronological order of their publishing date-time from YouTube for a given tag/search query in a paginated response.

# Basic Requirements:

- Server should call the YouTube API continuously in background (async) with some interval (say 10 seconds) for fetching the latest videos for a predefined search query and should store the data of videos (specifically these fields - Video title, description, publishing datetime, thumbnails URLs and any other fields you require) in a database with proper indexes.
- A GET API which returns the stored video data in a paginated response sorted in descending order of published datetime.
- A basic search API to search the stored videos using their title and description.
- Dockerize the project.
- It should be scalable and optimised.

# Bonus Points:

- Add support for supplying multiple API keys so that if quota is exhausted on one, it automatically uses the next available key.
- Make a dashboard to view the stored videos with filters and sorting options (optional)
- Optimise search api, so that it's able to search videos containing partial match for the search query in either video title or description.
    - Ex 1: A video with title *`How to make tea?`* should match for the search query `tea how`

# Instructions:

- You are free to choose any search query, for example: official, cricket, football etc. (choose something that has high frequency of video uploads)
- Try and keep your commit messages clean, and leave comments explaining what you are doing wherever it makes sense.
- Also try and use meaningful variable/function names, and maintain indentation and code style.
- Submission should have a `README` file containing instructions to run the server and test the API.


# Reference:

- YouTube data v3 API: [https://developers.google.com/youtube/v3/getting-started](https://developers.google.com/youtube/v3/getting-started)
- Search API reference: [https://developers.google.com/youtube/v3/docs/search/list](https://developers.google.com/youtube/v3/docs/search/list)
    - To fetch the latest videos you need to specify these: type=video, order=date, publishedAfter=<SOME_DATE_TIME>
    - Without publishedAfter, it will give you cached results which will be too old
    
 # API
 ```
 GET /health
 //should return {"status":"ok"} if server is operational
 
GET /search
// returns list of JSON object of truncated to a length of DEFAULT_RESULTS sorted in decreasing order by when they are published

GET /search?limit=26
// returns list of JSON object of truncated to a length of 26

GET /search?q=asian%20dog&limit=26
// returns list of JSON object of truncated after performing parital text search on both title and description truncated at a length of 26

GET /search?q=asian%20dog
// returns list of JSON object of truncated after performing parital text search on both title and description truncated at a length of DEFAULT_RESULTS
 
 ```
# DataBase Setup
We use a MongoDB database and setup 2 indexes .
- create a database name `Youtube` and a collection named `searchResult`
- create a search index with name as `default` .[follow this tutorial](https://www.mongodb.com/docs/atlas/atlas-search/create-index/)
```
{
  "mappings": {
    "dynamic": true
  }
}
```
- create a unique index by connecting via `mongosh` to avoid duplicate entries
```
Youtube.searchResult.createIndex( { title: 1, description: 1, channelTitle: 1 }, { unique: true } )
```
# How To Run Project using docker

-  pull docker image 
`docker pull karyamsettymartin/go-fampay:latest`

- set environment variables . (if you dont want the cron job to run on current instance set `CRON_INSTANCE=false`else set `CRON_INSTANCE=true`)
```
export API_KEY=AIzaSyCwCq28k0R_XEZIPtfoPPcD3XSVkZSymO0
export PORT=8000
export export DB_ADDRESS=mongodb+srv://martin:mishravikas@cluster0.k1p632w.mongodb.net/?retryWrites=true&w=majority
export CRON_INSTANCE=true
export DEFAULT_RESULTS=25
export MINING_INTERVAL_IN_MIN=10
export PUBLISHED_AFTER=020-09-01T01:59:53Z
export TOPIC=dogs`
```
- when scaling out horizontally keep in mind to only run cron job on only one instance as multiple cron jobs mining data will result in API key being blocked


- run docker image . and the server should be accessible from port 8000
```
 docker run -p 8000:8000 -host -it karyamsettymartin/go-fampay
```
# How To Run Project manually
- clone code from github
```
git clone https://github.com/ksnmartin/fampay.git
```
- set environment variables or use a `.env` file
```
export API_KEY=AIzaSyCwCq28k0R_XEZIPtfoPPcD3XSVkZSymO0
export PORT=8000
export export DB_ADDRESS=mongodb+srv://martin:mishravikas@cluster0.k1p632w.mongodb.net/?retryWrites=true&w=majority
export CRON_INSTANCE=true
export DEFAULT_RESULTS=25
export MINING_INTERVAL_IN_MIN=10
export PUBLISHED_AFTER=020-09-01T01:59:53Z
export TOPIC=dogs`
```
- build using compiler and run
```
go build -o run.exe
./run.exe
```
