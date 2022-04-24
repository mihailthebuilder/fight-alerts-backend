Feature: Scraper Service lambda scrapes sherdog.com for all the fight data

    Scenario: service goes to Sherdog and handles the full flow successfully

        Given fight records exist in the database
        When lambda is invoked
        Then the original fight records are deleted
        And newly-scraped fight records are inserted into the database