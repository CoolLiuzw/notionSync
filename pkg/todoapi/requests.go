package todoapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"notionsync/pkg/logger"
)

const urlPrefix = "https://graph.microsoft.com/beta/me/tasks/lists"

func (c *Client) CreateTaskList(name string) error {
	data := map[string]string{"displayName": name}
	req, err := NewJSONRequest(http.MethodPost, "", nil, data)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		return errors.New("response status code error")
	}

	return nil
}

func (c *Client) ListTaskLists() ([]TaskList, error) {
	req, err := NewJSONRequest(http.MethodGet, "", nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response status code error")
	}

	var list ListTaskListsResponse
	if err = json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}

	return list.TaskLists, nil
}

func (c *Client) GetTaskListByListName(listName string) ([]Task, error) {
	param := make(url.Values)
	param.Add("$filter", "contains(displayName,'"+listName+"')")
	req, err := NewRequest(http.MethodGet, "", nil, param, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response status code error")
	}

	var listTasks ListTasksResponse
	if err = json.NewDecoder(resp.Body).Decode(&listTasks); err != nil {
		return nil, err
	}

	return listTasks.Tasks, nil
}

func (c *Client) ListTask(taskListID string) ([]Task, error) {
	req, err := NewRequest(http.MethodGet, "/"+taskListID+"/tasks", nil, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response status code error")
	}

	var listTasks ListTasksResponse
	if err = json.NewDecoder(resp.Body).Decode(&listTasks); err != nil {
		return nil, err
	}

	return listTasks.Tasks, nil
}

func (c *Client) GetTaskDeltaLatest(taskListID string) (string, error) {
	param := make(url.Values)
	param.Add("$deltaToken", "latest")
	req, err := NewRequest(http.MethodGet, "/"+taskListID+"/tasks/delta", nil, param, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("response status code error")
	}
	// bodyByte, err := io.ReadAll(resp.Body)
	// if err != nil {
	//	return "", err
	// }
	// logger.Debugf("body:%v", string(bodyByte))
	// resp.Body = ioutil.NopCloser(bytes.NewReader(bodyByte))

	var jsonStruct struct {
		DeltaLink string `json:"@odata.deltaLink"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&jsonStruct); err != nil {
		return "", err
	}

	if len(jsonStruct.DeltaLink) == 0 {
		return "", errors.New("delta link is null")
	}

	return jsonStruct.DeltaLink, nil
}

func (c *Client) GetTaskDelta(taskListID string, inURL string) (*ListTasksResponse, error) {
	var uri string
	if inURL == "" {
		uri = "/" + taskListID + "/tasks/delta"
	} else {
		uri = strings.Replace(inURL, "https://graph.microsoft.com/beta/me/tasks/lists", "", -1)
	}
	req, err := NewRequest(http.MethodGet, uri, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		logger.Warnf("response status code: %v", resp.StatusCode)
		return nil, errors.New("response status code error")
	}

	var listTasks ListTasksResponse
	if err = json.NewDecoder(resp.Body).Decode(&listTasks); err != nil {
		return nil, err
	}

	return &listTasks, nil
}
