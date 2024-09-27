# file: features/event.feature

# http://localhost:8585/
# http://calendar-app:8585/

Feature: Calendar management
  In order to use calendar API
  As a responsible person
  I need to be able to add events

  Scenario: then user try to insert one event, one created event should be displayed by the system
    When I send POST request to "/v1/events" with payload:
        """
        {
            "title": "6th event",
            "startTime": "2024-12-30T10:00:00Z",
            "endTime": "2024-12-30T11:00:00Z",
            "description": "long description for 6th event",
            "userId": "4dde43c1-8381-41ab-8433-741d931081ff"
        }
        """
    Then the response code should be 200
    And the response payload should match json:
        """
        {
            "id": "bda84e66-321c-4fda-b4a8-a7bffae340df",
            "title": "6th event",
            "startTime": "2024-12-30T10:00:00Z",
            "endTime": "2024-12-30T11:00:00Z",
            "description": "long description for 6th event",
            "userId": "4dde43c1-8381-41ab-8433-741d931081ff"
        }
        """
