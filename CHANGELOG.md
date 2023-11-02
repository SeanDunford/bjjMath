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