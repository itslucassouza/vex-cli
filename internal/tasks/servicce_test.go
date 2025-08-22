package tasks

import "testing"

func TestAddTask(t *testing.T) {
	s := NewTaskService()
	task := s.AddTask("Comprar pão")

	if len(s.ListTasks()) != 1 {
		t.Fatal("esperava 1 tarefa")
	}

	if task.Title != "Comprar pão" {
		t.Errorf("esperava título 'Comprar pão', obteve '%s'", task.Title)
	}
}

func TestToggleTask(t *testing.T) {
	s := NewTaskService()
	task := s.AddTask("Fazer exercício")

	success := s.ToggleTask(task.ID)
	if !success {
		t.Fatal("esperava sucesso ao alternar tarefa")
	}

	if !s.ListTasks()[0].Done {
		t.Error("esperava tarefa marcada como feita")
	}
}

func TestRemoveTask(t *testing.T) {
	s := NewTaskService()
	task := s.AddTask("Estudar Go")
	s.RemoveTask(task.ID)

	if len(s.ListTasks()) != 0 {
		t.Error("esperava lista vazia após remoção")
	}
}
