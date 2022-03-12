package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

const jobsCounter, workersCounter = 100, 4

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	rand.Seed(time.Now().Unix())
	startTime := time.Now()

	jobs := make(chan int, jobsCounter)
	results := make(chan User, jobsCounter)

	for i := 0; i < workersCounter; i++ {
		go generateUsers(jobs, results)
	}

	for job := 0; job < jobsCounter; job++ {
		jobs <- job + 1
	}
	close(jobs)

	for job := 0; job < jobsCounter; job++ {
		go saveUserInfo(<-results)
	}

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func saveUserInfo(user User) {
	fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

	filename := fmt.Sprintf("users/uid%d.txt", user.id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(user.getActivityInfo())
	time.Sleep(time.Second)
}

func generateUsers(jobs <-chan int, results chan<- User) {
	for i := range jobs {
		user := User{
			id:    i,
			email: fmt.Sprintf("user%d@company.com", i),
			logs:  generateLogs(rand.Intn(10)),
		}
		fmt.Printf("generated user %d\n", i)
		time.Sleep(time.Millisecond * 100)

		results <- user
	}
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
