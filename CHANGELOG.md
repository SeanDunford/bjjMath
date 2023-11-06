# 0.0.6 -> NEXT
- Added a flag to run the scraper for one athelete 
- Tested all caching and parsing stages to ensure they work as expected 
- Going to further abstract the scraper so that shared code is reused


# 0.0.5 -> CURRENT
Implemented basic types for Athlete, Match, Record. The majority of athletes do not have any records on the page. Even People like Caio Terra have a completely empty bjj record which is suprising. Will be interesting to see if i can consume other datasources and start to try and build some type of historical record. Possibly even with crowd sourced submissions and adding other feeds like addcc wikipedia, ibjjf site and what not. De duplicating all of that data will be a problem in itself, but cross that bridge when we get there. 

# 0.0.4
Implemented scraping the athlete page. There's so much code duplication. Don't have a good abstracting/refactoring workflow with Go. There's so much hesitation to break things up becasue the compiler is so loud. 
I think I need to implement some types so there's less ambiguity in what things are. Now that the basic data types are defined. I'll abstract those as interfaces and replace all the []Athlete{} references to something liek []Athlete. and []AthleteRecord. After that i should find a way to persist them to a 
sqllite db using turso. https://turso.tech/libsql
Next Steps: 
  - Create Interfaces for basic data structures like the Athlete and their profile
  - Implement SQLLite to persist information 
  - Add some unit tests for the scraper 
  - Abstract scraper to be more generic 
  - Abstract csv methods to their own package 


# 0.0.3
Fixed limit, added a case for limit = -1 to use as no limit (make em say uhhhhhh)

# 0.0.2
Fully resolving urls and storing the mappings and basic athlete info in csv. Added sqllite as a dep but 
haven't implemented it yet. 

# 0.0.1
Basic functionality of getting a list of athletes and their profiles from BjjHEroes complete. 
Local caching also enabled for scraping and output to csv. 
Next Steps: 
  - Preserve a mapping of the unresolved to resolved urls in case we need later 
  - Implement SQLLite to persist information 
  - Implement a scraper for the individual athlete profiles and their records 
  - Create Interfaces for basic data structures like the Athlete and their profile
  - Test to see if i can just have a 'src' folder with packages inside. And require src/
