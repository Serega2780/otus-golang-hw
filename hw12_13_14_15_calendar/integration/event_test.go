package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/cucumber/godog"
)

type godogsResponseCtxKey struct{}

type eventTest struct {
	url string
}

type response struct {
	status int
	body   any
}

func InitializeScenario(godogCtx *godog.ScenarioContext) {
	test := &eventTest{url: "http://calendar-app:8585"}
	godogCtx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return context.WithValue(ctx, godogsResponseCtxKey{}, &response{status: 0, body: nil}), nil
	})
	godogCtx.Step(`^I send POST request to "([^"]*)" with payload:$`, test.iSendRequestToWithPayload)
	godogCtx.Step(`^the response code should be (\d+)$`, test.theResponseCodeShouldBe)
	godogCtx.Step(`^the response payload should match json:$`, test.theResponsePayloadShouldMatchJSON)

	godogCtx.Step(`^I send POST request to "([^"]*)" with payload:$`, test.iSendRequestToWithErrorPayload)
	godogCtx.Step(`^the response code should be (\d+)$`, test.theResponseCodeShouldBe)
	godogCtx.Step(`^the response payload should match json:$`, test.theResponsePayloadShouldMatchErrorJSON)

	godogCtx.Step(`^I send GET request to "([^"]*)" by week`, test.iSendRequestToFindByWeek)
	godogCtx.Step(`^the response code should be (\d+)$`, test.theResponseCodeShouldBe)
	godogCtx.Step(`^the response payload size should be (\d+)$`, test.theResponsePayloadSizeShouldMatch)

	godogCtx.Step(`^I send GET request to "([^"]*)" by month`, test.iSendRequestToFindByMonth)
	godogCtx.Step(`^the response code should be (\d+)$`, test.theResponseCodeShouldBe)
	godogCtx.Step(`^the response payload size should be (\d+)$`, test.theResponsePayloadSizeShouldMatch)

	godogCtx.Step(`^I send GET request to "([^"]*)" by day`, test.iSendRequestToFindByDay)
	godogCtx.Step(`^the response code should be (\d+)$`, test.theResponseCodeShouldBe)
	godogCtx.Step(`^the response payload size should be (\d+)$`, test.theResponsePayloadSizeShouldMatch)
}

func (et *eventTest) iSendRequestToFindByMonth(ctx context.Context, route string) (context.Context, error) {
	date := time.Now().Add(15 * 24 * time.Hour).Format(time.RFC3339)
	res, _ := http.Get(et.url + route + date) //nolint:noctx, bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	var events []model.Event
	_ = json.NewDecoder(res.Body).Decode(&events)

	actual := response{
		status: res.StatusCode,
		body:   len(events),
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, &actual), nil
}

func (et *eventTest) iSendRequestToFindByWeek(ctx context.Context, route string) (context.Context, error) {
	date := time.Now().Add(3 * 24 * time.Hour).Format(time.RFC3339)
	res, _ := http.Get(et.url + route + date) //nolint:noctx, bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	var events []model.Event
	_ = json.NewDecoder(res.Body).Decode(&events)
	actual := response{
		status: res.StatusCode,
		body:   len(events),
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, &actual), nil
}

func (et *eventTest) iSendRequestToFindByDay(ctx context.Context, route string) (context.Context, error) {
	date := time.Now().Format(time.RFC3339)
	res, _ := http.Get(et.url + route + date) //nolint:noctx, bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	var events []model.Event
	_ = json.NewDecoder(res.Body).Decode(&events)

	actual := response{
		status: res.StatusCode,
		body:   len(events),
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, &actual), nil
}

func (et *eventTest) theResponsePayloadSizeShouldMatch(ctx context.Context, expectedSze int) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	events := actualResp.body.(int)
	if events != expectedSze {
		return fmt.Errorf("expected response size does not match actual, %v vs. %v", expectedSze, events)
	}

	return nil
}

func (et *eventTest) theResponsePayloadShouldMatchErrorJSON(ctx context.Context, expectedBody *godog.DocString) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if actualResp.body.(string) != expectedBody.Content {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expectedBody, actualResp.body)
	}

	return nil
}

func (et *eventTest) iSendRequestToWithPayload(ctx context.Context, route string,
	payloadDoc *godog.DocString,
) (context.Context, error) {
	var reqBody []byte

	if payloadDoc != nil {
		payloadMap := &model.Event{}
		err := json.Unmarshal([]byte(payloadDoc.Content), payloadMap)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		reqBody, _ = json.Marshal(payloadMap)
	}
	res, _ := http.Post(et.url+route, "application/json", bytes.NewReader(reqBody)) //nolint:noctx, bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	createdEvent := &model.Event{}
	_ = json.NewDecoder(res.Body).Decode(createdEvent)

	actual := response{
		status: res.StatusCode,
		body:   *createdEvent,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, &actual), nil
}

func (et *eventTest) iSendRequestToWithErrorPayload(ctx context.Context, route string,
	payloadDoc *godog.DocString,
) (context.Context, error) {
	var reqBody []byte

	if payloadDoc != nil {
		payloadMap := &model.Event{}
		err := json.Unmarshal([]byte(payloadDoc.Content), payloadMap)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		reqBody, _ = json.Marshal(payloadMap)
	}
	res, _ := http.Post(et.url+route, "application/json", bytes.NewReader(reqBody)) //nolint:noctx, bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	bodyString := string(bodyBytes)
	fmt.Printf("%s\n", bodyString)

	actual := response{
		status: res.StatusCode,
		body:   bodyString,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, &actual), nil
}

func (et *eventTest) theResponseCodeShouldBe(ctx context.Context, expectedStatus int) error {
	resp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if expectedStatus != resp.status {
		if resp.status >= 400 {
			return fmt.Errorf("expected response code to be: %d, but actual is: %d,"+
				"response message: %s", expectedStatus, resp.status, resp.body)
		}
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", expectedStatus, resp.status)
	}

	return nil
}

func (et *eventTest) theResponsePayloadShouldMatchJSON(ctx context.Context, expectedBody *godog.DocString) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	expectedEvent := &model.Event{}

	err := json.Unmarshal([]byte(expectedBody.Content), expectedEvent)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	actualEvent := actualResp.body.(model.Event)
	expectedEvent.ID = actualEvent.ID

	if !reflect.DeepEqual(actualEvent, *expectedEvent) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expectedBody, actualResp.body)
	}

	return nil
}
