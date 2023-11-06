# How to Run 
- Clone repo 
- go mod tidy 
- mkdir output && mkdir output/html && output/csv


# References
- https://go.dev/doc/faq#git_https
- 


# Deps 
- http://go-colly.org/docs/introduction/install/
- 

# Future TO-DO 

## Beth stack for web site
- https://www.youtube.com/watch?v=cpzowDDJj24&t=0s

## Turso 
https://turso.tech/
Managed SQlLite as libSQL 
seems promising way to add db to app very quickly with flat file local dev and generous free tier 

## Add a compilation phase and build steps 
- Create output directories 
- add git version to output maybe

## How to work out belt rank for edges
A competitor that beat gordon at purple shouldn't be considered to beat gordon at black belt 

## How to work out scheduled jobs 
Refresh data daily/weekly/etc 

## How to figure out ruleset 
If a user beats gordon in EBI should that count for ADCC? What about points 

## [] Cache local files 
https://github.com/gocolly/colly/issues/193

## Expire local file cache
I need to be able to record when a file was cached locally and then use current date to validate/expire the cache date
I think for now once a month or week would really be fine. 

## [] REad local files 
https://github.com/gocolly/colly/blob/master/_examples/local_files/local_files.go

## [] Test new file against local file last modified date 
Update local cache and rescrape 

## Resolve urls 
Currently they are using weird ?p=# url and want the actual resolved urls 

## SqlLite

If you need to have sqlite first in your PATH, run:
  echo 'export PATH="/usr/local/opt/sqlite/bin:$PATH"' >> ~/.zshrc

For compilers to find sqlite you may need to set:
  export LDFLAGS="-L/usr/local/opt/sqlite/lib"
  export CPPFLAGS="-I/usr/local/opt/sqlite/include"

For pkg-config to find sqlite you may need to set:
  export PKG_CONFIG_PATH="/usr/local/opt/sqlite/lib/pkgconfig"

### graph db with sqllite
- https://github.com/actyx/ipfs-sqlite-block-store
- https://github.com/dpapathanasiou/simple-graph-go



# Gotchas
## Add protocol to scraping domain 
must send `https://www.bjjheroes.com` not just `www...` weirdly it just failed silently would have expected an error. 

## Adding local pacakages 
must run go mod init after 

## Best way to organize code 
I'd like to just have 

/src 
 - main.go
 /package1
  - package1.go
 /package2
   - package2.go
 ...


 but not sure entirely how to accomplish this. `main.go` should just setup the code with flags and then call
 into the src/ packages. 

 Go community pretty opinionated and have a standard
 https://github.com/golang-standards/project-layout

 ## Bad url 
 Sheliah	Lindsey		Magness BJJ
 has original url: https://www.bjjheroes.com/?p=13193
 which redirects to: https://www.bjjheroes.com/wp-content/uploads/2023/01/CarlosHenrique.jpg

it should be: https://www.bjjheroes.com/bjj-fighters/sheliah-lindsey


