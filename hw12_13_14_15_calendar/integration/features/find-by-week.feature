# file: features/event.feature

# http://localhost:8585/
# http://calendar-app:8585/

Feature: Calendar management
  In order to use calendar API
  As a responsible person
  I need to be able to find events by week

 Scenario: then user try to find events for a week, one event should be found
    When I send GET request to "/v1/events/week/" by week
    Then the response code should be 200
    And the response payload size should be 1
