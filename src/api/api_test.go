package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"testing"
)

// resetDatabase resets the database before running each test to ensure that the tests are independent of each other.
func resetDatabase() {
	println("resetting database")
	client := &http.Client{}
	url := fmt.Sprintf("http://localhost:%s/api/reset", getPort())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Printf("error during test setup: %s\n", err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error cleaning up test setup: %s\n", err)
		}
	}(res.Body)
}

// getPort returns the port number specified in .env.
func getPort() string {
	// utils.GetEnvVariable() is not used here because it won't be able to find the .env file
	// without additional configuration.
	// Instead, a relative path is the simplest solution to this problem.
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Printf("error loading .env: %s\n", err)
	}
	return os.Getenv("PORT")
}

// setRequest encapsulates the boilerplate code for creating a request with an optional JSON body.
func sendRequest(t *testing.T, route string, method string, json string) *http.Response {
	url := fmt.Sprintf("http://localhost:%s/%s", getPort(), route)
	var reader *bytes.Reader
	reader = bytes.NewReader([]byte(json))
	req, err := http.NewRequest(method, url, reader)
	if json != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if err != nil {
		t.Errorf("client: could not create request: %s\n", err)
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("client: could not create request: %s\n", err)
	}
	return res
}

// TestRegister tests that students can be registered to a teacher.
func TestRegister(t *testing.T) {
	resetDatabase()
	res := sendRequest(t, "api/register", "POST", `
		{
		  "teacher": "teacherken@gmail.com",
		  "students": [
			"studentjon@gmail.com",
			"studenthon@gmail.com"
		  ]
		}
	`)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

// TestCommonStudents tests that common students can be retrieved.
func TestCommonStudents(t *testing.T) {
	resetDatabase()
	res := sendRequest(t, "api/register", "POST", `
		{
		  "teacher": "teacherken@gmail.com",
		  "students": [
			"commonstudent1@gmail.com",
			"commonstudent2@gmail.com",
			"student_only_under_teacher_ken@gmail.com"
		  ]
		}
	`)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	res = sendRequest(t, "api/register", "POST", `
		{
		  "teacher": "teacherjoe@gmail.com",
		  "students": [
			"commonstudent1@gmail.com",
			"commonstudent2@gmail.com"
		  ]
		}
	`)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	res = sendRequest(t, "api/commonstudents?teacher=teacherken%40gmail.com", "GET", "")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var students Students
	err := json.NewDecoder(res.Body).Decode(&students)
	if err != nil {
		t.Errorf("could not decode response: %s\n", err)
	}
	assert.Contains(t, students.Students, "commonstudent1@gmail.com")
	assert.Contains(t, students.Students, "commonstudent2@gmail.com")
	assert.Contains(t, students.Students, "student_only_under_teacher_ken@gmail.com")
	assert.Equal(t, 3, len(students.Students))
}

func TestSuspendNonExistentStudent(t *testing.T) {
	resetDatabase()
	res := sendRequest(t, "api/suspend", "POST", `
		{
		  "student": "studentmary@gmail.com"
		}
	`)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestSuspendStudent(t *testing.T) {
	resetDatabase()
	res := sendRequest(t, "api/register", "POST", `
		{
		  "teacher": "teacherken@gmail.com",
		  "students": [
			"studentmary@gmail.com"
		  ]
		}
	`)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	res = sendRequest(t, "api/suspend", "POST", `
		{
		  "student": "studentmary@gmail.com"
		}
	`)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestRetrieveForNotifications(t *testing.T) {
	resetDatabase()
	res := sendRequest(t, "api/register", "POST", `
		{
		  "teacher": "teacherken@gmail.com",
		  "students": [
			"studentbob@gmail.com"
		  ]
		}
	`)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	res = sendRequest(t, "api/register", "POST", `
	{
	  "teacher": "teacherjen@gmail.com",
	  "students": [
		"studentagnes@gmail.com",
		"studentmiche@gmail.com"
	  ]
	}
	`)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	res = sendRequest(t, "api/retrievefornotifications", "POST", `
		{
		"teacher": "teacherken@gmail.com",
		"notification": "Hey everybody"
		}
	`)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var recipients Recipients
	err := json.NewDecoder(res.Body).Decode(&recipients)
	if err != nil {
		t.Errorf("could not decode response: %s\n", err)
	}
	assert.Contains(t, recipients.Recipients, "studentbob@gmail.com")
	assert.Equal(t, 1, len(recipients.Recipients))
	res = sendRequest(t, "api/retrievefornotifications", "POST", `
		{
		  "teacher": "teacherken@gmail.com",
		  "notification": "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com"
		}
	`)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(&recipients)
	if err != nil {
		t.Errorf("could not decode response: %s\n", err)
	}
	assert.Contains(t, recipients.Recipients, "studentbob@gmail.com")
	assert.Contains(t, recipients.Recipients, "studentagnes@gmail.com")
	assert.Contains(t, recipients.Recipients, "studentmiche@gmail.com")
	assert.Equal(t, 3, len(recipients.Recipients))
}
