package notion

import (
	"context"
	"time"

	"notionsync/pkg/notionapi"

	"github.com/pkg/errors"
)

var (
	_true  = true
	_false = false
)

type API interface {
	AddTask(title, todoID, importance, displayName string) error
	AddTaskWithScheduleTime(title, todoID, importance, displayName, scheduleTimeStr string) error
	CompleteTask(title string) error
	ExistTaskFromTodoID(todoID string) (bool, error)
	UpdateTaskInfo(todoID, title, status, importance, dueDateTime, taskListName string, completedDateTime time.Time, deleted bool) error
}

type options struct {
	apiSecret  string
	databaseID string
}

type notion struct {
	ctx    context.Context
	cancel context.CancelFunc
	client *notionapi.Client
	option options
	pageID string
}

func New(apiSecret, databaseID string) API {
	option := options{
		apiSecret:  apiSecret,
		databaseID: databaseID,
	}

	ctx, cancel := context.WithCancel(context.TODO())
	return &notion{
		ctx:    ctx,
		cancel: cancel,
		client: notionapi.NewClient(apiSecret),
		option: option,
	}
}

func (n *notion) UpdateTaskInfo(todoID string, title string, status string, importance string, dueDateTime string, taskListName string, completedDateTime time.Time, deleted bool) error {
	queryDatabase, err := n.client.QueryDatabase(n.ctx, n.option.databaseID, &notionapi.DatabaseQuery{
		Filter: &notionapi.DatabaseQueryFilter{
			And: []notionapi.DatabaseQueryFilter{
				{
					Property: "TodoID",
					Text: &notionapi.TextDatabaseQueryFilter{
						Equals: todoID,
					},
				},
			},
		},
	})
	if err != nil {
		return errors.WithMessagef(err, "database query failed:%v:%v", n.option.databaseID, todoID)
	}

	if queryDatabase.HasMore {
		return errors.Errorf("query database id: %v, filter title: %v has more", n.option.databaseID, todoID)
	}

	if len(queryDatabase.Results) == 0 {
		return errors.Errorf("query database id: %v, filter title: %v not found", n.option.databaseID, todoID)
	}

	page := queryDatabase.Results[0]
	var databasePageProperties = make(notionapi.DatabasePageProperties)

	if len(status) > 0 {
		done := false
		if status == "completed" {
			done = true
		}
		databasePageProperties["Done"] = notionapi.DatabasePageProperty{
			Checkbox: &done,
		}
	}
	databasePageProperties["Deleted"] = notionapi.DatabasePageProperty{
		Checkbox: &deleted,
	}

	if len(dueDateTime) > 0 {
		const timeLayout = "2006-01-02T15:04:05.0000000"
		parse, err := time.Parse(timeLayout, dueDateTime)
		if err != nil {
			panic(err)
		}
		scheduleTime := parse.Add(time.Hour * 24)
		databasePageProperties["Scheduled Time"] = notionapi.DatabasePageProperty{
			Date: &notionapi.Date{
				Start:    notionapi.NewDateTime(scheduleTime, false),
				End:      nil,
				TimeZone: nil,
			}}
	}

	if !completedDateTime.IsZero() {
		databasePageProperties["Completion time"] = notionapi.DatabasePageProperty{
			Date: &notionapi.Date{
				Start:    notionapi.NewDateTime(completedDateTime, false),
				End:      nil,
				TimeZone: nil,
			}}
	}

	if len(title) > 0 {
		databasePageProperties["Task"] = notionapi.DatabasePageProperty{
			Title: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: title,
					},
				},
			},
		}
	}

	if len(importance) > 0 {
		var s notionapi.SelectOptions
		if importance == "normal" {
			s.Name = "P2"
		} else if importance == "high" {
			s.Name = "P0 🔥"
		} else {
			s.Name = "P2"
		}
		databasePageProperties["Priority"] = notionapi.DatabasePageProperty{
			Select: &s,
		}
	}

	if len(taskListName) > 0 {
		databasePageProperties["Task List Name"] = notionapi.DatabasePageProperty{
			RichText: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: taskListName,
					},
				},
			},
		}
	}

	_, err = n.client.UpdatePage(n.ctx, page.ID, notionapi.UpdatePageParams{
		DatabasePageProperties: &databasePageProperties,
	})
	if err != nil {
		return errors.WithMessagef(err, "update database %v, page %v failed", n.option.databaseID, page.ID)
	}
	return nil
}

func (n *notion) ExistTaskFromTodoID(todoID string) (bool, error) {
	database, err := n.client.QueryDatabase(n.ctx, n.option.databaseID, &notionapi.DatabaseQuery{
		Filter: &notionapi.DatabaseQueryFilter{
			And: []notionapi.DatabaseQueryFilter{
				{
					Property: "TodoID",
					Text: &notionapi.TextDatabaseQueryFilter{
						Equals: todoID,
					},
				},
			},
		},
	})
	if err != nil {
		return false, errors.WithMessagef(err, "exist database query failed:%v:%v", n.option.databaseID, todoID)
	}

	if len(database.Results) == 0 {
		return false, nil
	}

	return true, nil
}

func (n *notion) AddTask(title, todoID, importance, displayName string) error {
	return n.addTask(title, todoID, importance, displayName, "")
}

func (n *notion) AddTaskWithScheduleTime(title, todoID, importance, displayName, scheduleTimeStr string) error {
	return n.addTask(title, todoID, importance, displayName, scheduleTimeStr)
}

func (n *notion) addTask(title, todoID, importance, displayName, scheduleTime string) error {
	database, err := n.client.FindDatabaseByID(n.ctx, n.option.databaseID)
	if err != nil {
		return errors.WithMessagef(err, "add task database id: %v failed", n.option.databaseID)
	}

	var databasePageProperties = make(notionapi.DatabasePageProperties)

	if len(importance) > 0 {
		var s notionapi.SelectOptions
		if importance == "normal" {
			s.Name = "P2"
		} else if importance == "high" {
			s.Name = "P0 🔥"
		} else {
			s.Name = "P2"
		}
		databasePageProperties["Priority"] = notionapi.DatabasePageProperty{
			Select: &s,
		}
	}

	databasePageProperties["Task"] = notionapi.DatabasePageProperty{
		Title: []notionapi.RichText{
			{
				Text: &notionapi.Text{
					Content: title,
				},
			},
		},
	}
	databasePageProperties["TodoID"] = notionapi.DatabasePageProperty{
		RichText: []notionapi.RichText{
			{
				Text: &notionapi.Text{
					Content: todoID,
				},
			},
		},
	}
	databasePageProperties["Task List Name"] = notionapi.DatabasePageProperty{
		RichText: []notionapi.RichText{
			{
				Text: &notionapi.Text{
					Content: displayName,
				},
			},
		},
	}

	if len(scheduleTime) > 0 {
		const timeLayout = "2006-01-02T15:04:05.0000000"
		parse, err := time.Parse(timeLayout, scheduleTime)
		if err != nil {
			panic(err)
		}
		scheduleTime := parse.Add(time.Hour * 24)
		databasePageProperties["Scheduled Time"] = notionapi.DatabasePageProperty{
			Date: &notionapi.Date{
				Start:    notionapi.NewDateTime(scheduleTime, false),
				End:      nil,
				TimeZone: nil,
			}}
	}

	_, err = n.client.CreatePage(
		n.ctx,
		notionapi.CreatePageParams{
			ParentType:             notionapi.ParentTypeDatabase,
			ParentID:               database.ID,
			DatabasePageProperties: &databasePageProperties,
		},
	)
	if err != nil {
		return errors.WithMessagef(err, "database id: %v, create page failed", n.option.databaseID)
	}

	return nil
}

func (n *notion) CompleteTask(title string) error {
	queryDatabase, err := n.client.QueryDatabase(n.ctx, n.option.databaseID, &notionapi.DatabaseQuery{
		Filter: &notionapi.DatabaseQueryFilter{
			And: []notionapi.DatabaseQueryFilter{
				{
					Property: "Done",
					Checkbox: &notionapi.CheckboxDatabaseQueryFilter{
						Equals: &_false,
					},
				},
				{
					Property: "Task",
					Text: &notionapi.TextDatabaseQueryFilter{
						Equals: title,
					},
				},
			},
		},
	})
	if err != nil {
		return errors.WithMessagef(err, "query database id: %v failed", n.option.databaseID)
	}

	if queryDatabase.HasMore {
		return errors.Errorf("query database id: %v, filter title: %v has more", n.option.databaseID, title)
	}

	if len(queryDatabase.Results) == 0 {
		return errors.Errorf("query database id: %v, filter title: %v not found", n.option.databaseID, title)
	}

	page := queryDatabase.Results[0]
	_, err = n.client.UpdatePage(n.ctx, page.ID, notionapi.UpdatePageParams{
		DatabasePageProperties: &notionapi.DatabasePageProperties{
			"Done": notionapi.DatabasePageProperty{
				Checkbox: &_true,
			},
			"Completion time": notionapi.DatabasePageProperty{
				Date: &notionapi.Date{
					Start:    notionapi.NewDateTime(time.Now(), true),
					End:      nil,
					TimeZone: nil,
				},
			},
		},
	})
	if err != nil {
		return errors.WithMessagef(err, "update database %v, page %v failed", n.option.databaseID, page.ID)
	}

	return nil
}
