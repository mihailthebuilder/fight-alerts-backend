Feature: Scraper Service scrapes sherdog.com for all the fight data

    Scenario: Scraper Service goes to Sherdog and handles the full flow successfully

        Given Sherdog is available for access
        When the service is invoked
        Then all the fight data is logged