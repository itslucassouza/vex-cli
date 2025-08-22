// task_service.go
package tasks

type TaskService struct {
	tasks  []Task
	nextID int
}

func NewTaskService() *TaskService {
	tasks, _ := LoadTasks()
	nextID := 1
	for _, t := range tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}

	return &TaskService{
		tasks:  tasks,
		nextID: nextID,
	}
}

func (s *TaskService) AddTask(title string) Task {
	task := Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks = append(s.tasks, task)
	s.nextID++
	_ = SaveTasks(s.tasks)
	return task
}

func (s *TaskService) ToggleTask(id int) bool {
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks[i].Done = !s.tasks[i].Done
			_ = SaveTasks(s.tasks)
			return true
		}
	}
	return false
}

func (s *TaskService) RemoveTask(id int) bool {
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			_ = SaveTasks(s.tasks)
			return true
		}
	}
	return false
}

func (s *TaskService) EditTask(id int, newTitle string) bool {
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks[i].Title = newTitle
			_ = SaveTasks(s.tasks)
			return true
		}
	}
	return false
}

func (s *TaskService) ListTasks() []Task {
	return s.tasks
}
