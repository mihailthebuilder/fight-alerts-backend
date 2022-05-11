Feature: Scraper Service lambda scrapes sherdog.com for all the fight data

    Scenario: service goes to Sherdog and handles the full flow successfully

        Given some fight records exist in the database
        And a trigger for the notification service has been set in eventbridge

        When the scraper lambda is invoked

        Then the original fight records are deleted
        And newly-scraped fight records are inserted into the database
        And the notification service trigger is updated