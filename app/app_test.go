package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"stacew/todos/dataModel"

	"github.com/stretchr/testify/assert"
)

func TestTodos(t *testing.T) {
	getSessionID = func(r *http.Request) string {
		return "testsessionID"
	}

	assert := assert.New(t)

	testFilePath := "./test.db"
	os.Remove(testFilePath)
	ah := MakeNewHandler(testFilePath)
	defer ah.Close()

	ts := httptest.NewServer(ah)
	defer ts.Close()

	resp, err := http.PostForm(ts.URL+"/todoH", url.Values{"name": {"Test todo"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	var todo dataModel.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo")
	id1 := todo.ID

	resp, err = http.PostForm(ts.URL+"/todoH", url.Values{"name": {"Test todo2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	var todo2 dataModel.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo2)
	assert.NoError(err)
	assert.Equal(todo2.Name, "Test todo2")
	id2 := todo2.ID

	resp, err = http.Get(ts.URL + "/todoH")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*dataModel.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)
	for _, t := range todos {
		if t.ID == id1 {
			assert.Equal("Test todo", t.Name)
		} else if t.ID == id2 {
			assert.Equal("Test todo2", t.Name)
		} else {
			assert.Error(fmt.Errorf("test ID should be id1 or id2"))
		}
	}

	resp, err = http.Get(ts.URL + "/complete-todoH/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(ts.URL + "/todoH")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*dataModel.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)
	for _, t := range todos {
		if t.ID == id1 {
			assert.True(t.Completed)
		}
	}

	req, _ := http.NewRequest("DELETE", ts.URL+"/todoH/"+strconv.Itoa(id1), nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(ts.URL + "/todoH")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*dataModel.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)
	for _, t := range todos {
		if t.ID == id2 {
			assert.Equal(t.ID, id2)
		}
	}

}
