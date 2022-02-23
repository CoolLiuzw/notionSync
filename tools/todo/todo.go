package todo

import (
	"log"
	"math/rand"
	"time"

	"notionsync/pkg/logger"
	"notionsync/pkg/todoapi"
	"notionsync/tools/notion"
)

type API interface {
	UpdateNotionAllToDo() error
}

type todo struct {
	client *todoapi.Client
	notion notion.API
}

func New(clientID, clientSecret string, notionAPI notion.API) (API, error) {
	client, err := todoapi.NewClient(clientID, clientSecret)
	if err != nil {
		return nil, err
	}

	return &todo{
		client: client,
		notion: notionAPI,
	}, nil
}

func GetToken(clientID, clientSecret string) {
	token, err := todoapi.GetToken(clientID, clientSecret)
	if err != nil {
		log.Fatal("Can't get token. Check your network and try again")
		return
	}

	log.Printf("%v", token.Expiry.String())
	saveTokenErr := todoapi.SaveToken([]byte(token.RefreshToken))
	if saveTokenErr != nil {
		log.Fatal("Can't save authenticating information. Try again")
		return
	}
}

func getTaskDeltaUrl(tasks *todoapi.ListTasksResponse) (isDeltaLink bool, url string) {
	if len(tasks.OdataDeltaLink) > 0 {
		return true, tasks.OdataDeltaLink
	}
	return false, tasks.OdataNextLink
}

func (t *todo) notionDeleteTask(tasksID string) {
	if err := t.notion.UpdateTaskInfo(tasksID, "", "", "", "", "", time.Time{}, true); err != nil {
		logger.Warnf("deleted task id failed: %v", err)
	}
}

func (t *todo) notionUpdateTaskInfo(task todoapi.Task, displayName string) {
	logger.Debugf("task update >>>> : [%v]", task.DisplayName)
	err := t.notion.UpdateTaskInfo(task.Id, task.DisplayName, task.Status, task.Importance, task.DueDateTime.DateTime, displayName, task.CompletedDateTime, false)
	if err != nil {
		logger.Warnf("notion update task info: %v failed, displayName: %v", err, displayName)
	}
}

func (t *todo) notinAddTaskInfo(task todoapi.Task, displayName string) {
	logger.Debugf("task create >>>> : [%v]", task.DisplayName)
	if len(task.DueDateTime.DateTime) == 0 {
		err := t.notion.AddTask(task.DisplayName, task.Id, task.Importance, displayName)
		if err != nil {
			logger.Warnf("notion add task: %v failed, displayName: %v", err, displayName)
		}
		return
	}

	if err := t.notion.AddTaskWithScheduleTime(task.DisplayName, task.Id,
		task.Importance, displayName, task.DueDateTime.DateTime); err != nil {
		logger.Warnf("notion add task: %v failed", err)
	}
}

func (t *todo) deltaLoop(taskListID, displayName string) {
	var tasks = &todoapi.ListTasksResponse{}

	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)

	logger.Debugf(taskListID + "::::" + displayName + "loop will start")

	for {
		deltaLink, url := getTaskDeltaUrl(tasks)
		respTask, err := t.client.GetTaskDelta(taskListID, url)
		if err != nil {
			logger.Warnf("get task delta: %v failed, displayName: %v", err, displayName)
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			continue
		}
		tasks = respTask

		if !deltaLink {
			logger.Debugf("not delta link: %v, will next", displayName)
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			continue
		}

		for _, task := range tasks.Tasks {
			if task.Removed.Reason == "deleted" {
				t.notionDeleteTask(task.Id)
				continue
			}

			if len(task.DisplayName) == 0 {
				logger.Warnf("task displayName is empty")
				continue
			}

			exist, err := t.notion.ExistTaskFromTodoID(task.Id)
			if err != nil {
				logger.Warnf("notion exist task from todo id failed: %v", err)
				continue
			}

			if exist {
				t.notionUpdateTaskInfo(task, displayName)
			} else {
				t.notinAddTaskInfo(task, displayName)
			}
		}

		var randSec int
		for {
			randSec = rand.Intn(60)
			if randSec > 30 {
				break
			}
		}
		logger.Debugf("time update now: %v, random second: %vs", displayName, randSec)
		time.Sleep(time.Duration(randSec) * time.Second)
	}
}

func (t *todo) UpdateNotionAllToDo() error {
	listTaskLists, err := t.client.ListTaskLists()
	if err != nil {
		return err
	}

	logger.Debugf("list len: %v", len(listTaskLists))
	for _, taskLists := range listTaskLists {
		go t.deltaLoop(taskLists.Id, taskLists.DisplayName)
	}

	select {}
}
