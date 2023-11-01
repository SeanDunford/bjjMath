# References
- https://go.dev/doc/faq#git_https
- 


# Deps 
- http://go-colly.org/docs/introduction/install/
- 

# Future TO-DO 

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


## In memory graph db 
- https://github.com/actyx/ipfs-sqlite-block-store
- https://github.com/dpapathanasiou/simple-graph-go



# Gotchas
## Add protocol to scraping domain 
must send `https://www.bjjheroes.com` not just `www...` weirdly it just failed silently would have expected an error. 