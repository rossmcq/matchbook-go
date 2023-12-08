# Golang Matchbook Odds Scraper

## Aim
Have an API that you can pass an event_id then it will start recording the Match Odds to a database to be analysed post or during game. 

NOTE: I have a very similar app writen in python `matchbook-python`

## Requirements
To be able to fetch the matchbook token you'll need env variables `MATCHBOOK_USER` and `MATCHBOOK_PW` populated with valid credentials to www.matchbook.com.

For Database connection you'll need to run the init.sql script and create a user that has access to the `matchbook` database. You'll need t ostore that database users credentials in env variables `POSTGRES_USER` and `POSTGRES_PASSWORD`
