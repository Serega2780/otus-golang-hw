# file: features/event.feature

# http://localhost:8585/
# http://calendar-app:8585/

Feature: Calendar management
  In order to use calendar API
  As a responsible person
  I need to be able to add events with error

  Scenario: then user try to insert an event with wrong date from the past, an error should appear
    When I send POST request to "/v1/events" with payload:
        """
        {
            "title": "7th event",
            "startTime": "2023-12-30T10:00:00Z",
            "endTime": "2023-12-30T11:00:00Z",
            "description": "long description for 7th event",
            "userId": "4dde43c1-8381-41ab-8433-741d931081ff"
        }
        """
    Then the response code should be 400
    And the response payload should match json:
        """
        {
            "error": "mandatory fields are missing or wrong error"
        }
        """