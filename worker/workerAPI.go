package worker

import (
	"encoding/json"
	"net/http"
)

func handleGetStatus(w http.ResponseWriter, r *http.Request, status Status, config WorkerConfig) {
	oauthKeyClient := r.Header.Get("Authorization")
	if oauthKeyClient != config.OAuthToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

/*

func handletaskMessage(w http.ResponseWriter, r *http.Request) {
	oauthKey := r.Header.Get("Authorization")
	if oauthKey != oauthToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("task message (ID: %s): %s\n", message.ID, message.Module)

	// Create a new Task
	task := &Task{
		ID:          message.ID,
		Module:      message.Module,
		Arguments:   message.Arguments,
		CallbackURL: message.CallbackURL,
		Status:      "Pending",
		Goroutine:   &sync.WaitGroup{},
	}

	// Add the task to the list
	taskListMu.Lock()
	taskList[message.ID] = task
	taskListMu.Unlock()

	// Start a new goroutine for the task
	task.Goroutine.Add(1)
	go processTask(message, task)

	// Respond immediately without waiting for the task to complete
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, message.ID)
}

func processTask(message Message, task *Task) {
	defer func() {
		task.Goroutine.Done()
		// Release a slot in the semaphore when the task is done
		<-semaphoreCh
	}()

	// Acquire a slot from the semaphore
	semaphoreCh <- struct{}{}

	//Set task status
	task.Status = "Working"

	workMutex.Lock()
	isWorking = true
	workMutex.Unlock()

	// Process the module in the task
	m, err := processModule(message.Module, message.Arguments)
	if err != 0 {
		fmt.Printf("Failed to run module")
	}

	workMutex.Lock()
	isWorking = false
	workMutex.Unlock()

	//Set task status
	task.Status = "Done"

	// Remove the task from the list
	//taskListMu.Lock()
	//delete(taskList, task.ID)
	//taskListMu.Unlock()

	// Save the output in the task
	task.Result = m

	payload, _ := json.Marshal(task)
	callbackTaskMessage(task.CallbackURL, "application/json", payload)
}

func callbackTaskMessage(url, contentType string, payload []byte) {
	// Create a new request with the POST method and the payload
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add custom headers, including the OAuth header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", oauthToken)

	// Create an HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.Status)
	// Handle the response body as needed

}

func handleGetStatus(w http.ResponseWriter, r *http.Request) {
	oauthKey := r.Header.Get("Authorization")
	if oauthKey != oauthToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workMutex.Lock()
	defer workMutex.Unlock()

	status := Status{
		IsWorking: isWorking,
		MessageID: messageID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func handleGetTasks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	status := r.URL.Query().Get("status")

	taskListMu.Lock()
	defer taskListMu.Unlock()

	var filteredTasks map[string]Task

	// Filter tasks by status if the status parameter is provided
	if status != "" {
		filteredTasks = make(map[string]Task)
		for id, task := range taskList {
			if status == task.Status {
				filteredTasks[id] = *task
			}
		}
	} else {
		filteredTasks = make(map[string]Task)
		for id, task := range taskList {
			filteredTasks[id] = *task
		}
	}

	responseJSON, err := json.Marshal(filteredTasks)
	if err != nil {
		http.Error(w, "Error encoding tasks to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	taskListMu.Lock()
	defer taskListMu.Unlock()

	task, exists := taskList[taskID]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Error encoding task info to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func processModule(module string, arguments []string) (string, int) {
	switch module {
	case "work1":
		workAndNotify(1, messageID)
		return "Task scheduled for work with an unknown duration", 0
	case "module1":
		return module1(arguments)
	case "module2":
		return module2(arguments)
	case "workList":
		if len(arguments) > 0 {
			// Simulate work with an unknown duration
			workDuration := getRandomDuration()
			time.Sleep(workDuration)
			return stringList(arguments), 0
		}
		return "", 1
	default:
		return "Unknown task", 0
	}
}
*/
