package todoapi

import (
	"time"

	oauth "golang.org/x/oauth2"
)

type TokenResponse struct {
	TokenValue *oauth.Token
}

type TaskList struct {
	OdataType         string `json:"@odata.type"`
	OdataEtag         string `json:"@odata.etag"`
	WellKnownListName string `json:"wellKnownListName"`
	DisplayName       string `json:"displayName"`
	Id                string `json:"id"`
}

type ListTaskListsResponse struct {
	DataContext string     `json:"@odata.context"`
	TaskLists   []TaskList `json:"value"`
}

type Task struct {
	OdataType            string    `json:"@odata.type"`
	OdataEtag            string    `json:"@odata.etag"`
	CompletedDateTime    time.Time `json:"completedDateTime"`
	Importance           string    `json:"importance"`
	Status               string    `json:"status"`
	DisplayName          string    `json:"displayName"`
	CreatedDateTime      time.Time `json:"createdDateTime"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
	Id                   string    `json:"id"`
	Body                 struct {
		Content     string `json:"content"`
		ContentType string `json:"contentType"`
	} `json:"body"`
	DueDateTime struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"dueDateTime"`
	StartDateTime struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"startDateTime"`
	ParentList struct {
		Id string `json:"id"`
	} `json:"parentList"`
	Removed struct {
		Reason string `json:"reason"`
	} `json:"@removed"`
}

type ListTasksResponse struct {
	OdataContext   string `json:"@odata.context"`
	OdataNextLink  string `json:"@odata.nextLink"`
	OdataDeltaLink string `json:"@odata.deltaLink"`
	Tasks          []Task `json:"value"`
}

type TaskBody struct {
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
}

type DateStruct struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timeZone"`
}

func (response TokenResponse) Token() (*oauth.Token, error) {
	return response.TokenValue, nil
}
