# file: features/event.feature

# http://localhost:8585/
# http://calendar-app:8585/

Feature: Calendar management
  In order to use calendar API
  As a responsible person
  I need to be able to find events by day

  Scenario: then user try to find events for today, three events should be found
    When I send GET request to "/v1/events/day/" by day
    Then the response code should be 200
    And the response payload size should be 3

