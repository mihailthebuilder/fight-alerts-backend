Feature: Scraper Service lambda scrapes sherdog.com for all the fight data

    Scenario: service goes to Sherdog and handles the full flow successfully
        Given Sherdog is available for access
        When lambda is invoked
        Then all the fight data is returned