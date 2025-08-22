package googletasks

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/cli/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	tasks "google.golang.org/api/tasks/v1"
)

// Retorna o cliente autenticado
func getClient(config *oauth2.Config) *http.Client {
	tokenFile := "token.json"

	tok, err := tokenFromFile(tokenFile)
	if err != nil || !tok.Valid() { // se token não existe ou inválido
		tok = getTokenFromWebAuto(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Lê token do arquivo
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Salva token no arquivo
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Salvando token em %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("Erro ao salvar token: %v\n", err)
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Fluxo automático de login via navegador e callback local
func getTokenFromWebAuto(config *oauth2.Config) *oauth2.Token {
	ctx := context.Background()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)

	server := &http.Server{Addr: ":9999"}

	// create a channel to receive the authorization code
	codeChan := make(chan string)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// get the OAuth authorization URL
	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above
	fmt.Printf("Your browser has been opened to visit::\n%s\n", url)

	// open user's browser to login page
	if err := browser.OpenURL(url); err != nil {
		panic(fmt.Errorf("failed to open browser for authentication %s", err.Error()))
	}

	// wait for the authorization code to be received
	code := <-codeChan

	// exchange the authorization code for an access token
	_, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Failed to exchange authorization code for token: %v", err)
	}

	// shut down the HTTP server
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Failed to shut down server: %v", err)
	}

	return nil

}

// Abre o navegador padrão
func openBrowser(url string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		return fmt.Errorf("sistema não suportado")
	}
	return exec.Command(cmd, args...).Start()
}

// Listar tasks
func ListTasks() ([]*tasks.Task, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("erro ao ler credentials.json: %v", err)
	}

	config, err := google.ConfigFromJSON(b, tasks.TasksScope)
	if err != nil {
		return nil, fmt.Errorf("erro no ConfigFromJSON: %v", err)
	}

	client := getClient(config)

	srv, err := tasks.New(client)
	if err != nil {
		return nil, fmt.Errorf("não foi possível criar o serviço: %v", err)
	}

	taskLists, err := srv.Tasklists.List().Do()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar tasklists: %v", err)
	}

	if len(taskLists.Items) == 0 {
		fmt.Println("Nenhuma lista de tarefas encontrada.")
		return nil, nil
	}

	var allTasks []*tasks.Task
	for _, tl := range taskLists.Items {
		tasksRes, err := srv.Tasks.List(tl.Id).Do()
		if err != nil {
			return nil, fmt.Errorf("erro ao listar tasks: %v", err)
		}
		allTasks = append(allTasks, tasksRes.Items...)
	}

	return allTasks, nil
}

// Adicionar task
func AddTask(taskListID, title string) error {
	ctx := context.Background()

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return fmt.Errorf("erro ao ler credentials.json: %v", err)
	}

	config, err := google.ConfigFromJSON(b, tasks.TasksScope)
	if err != nil {
		return fmt.Errorf("erro no ConfigFromJSON: %v", err)
	}

	client := getClient(config)

	srv, err := tasks.New(client)
	if err != nil {
		return fmt.Errorf("não foi possível criar o serviço: %v", err)
	}

	task := &tasks.Task{
		Title: title,
	}

	_, err = srv.Tasks.Insert(taskListID, task).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("erro ao adicionar tarefa: %v", err)
	}

	return nil
}

func GetDefaultTaskListID() (string, error) {
	ctx := context.Background()

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return "", fmt.Errorf("erro ao ler credentials.json: %v", err)
	}

	config, err := google.ConfigFromJSON(b, tasks.TasksScope)
	if err != nil {
		return "", fmt.Errorf("erro no ConfigFromJSON: %v", err)
	}

	client := getClient(config)

	srv, err := tasks.New(client)
	if err != nil {
		return "", fmt.Errorf("não foi possível criar o serviço: %v", err)
	}

	taskLists, err := srv.Tasklists.List().Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao listar tasklists: %v", err)
	}

	if len(taskLists.Items) == 0 {
		return "", fmt.Errorf("nenhuma lista de tarefas encontrada")
	}

	return taskLists.Items[0].Id, nil
}

func RemoveTask(taskListID, id string) error {
	ctx := context.Background()

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de credenciais: %v", err)
	}

	config, err := google.ConfigFromJSON(b, tasks.TasksScope)
	if err != nil {
		return fmt.Errorf("erro ao criar config do JSON: %v", err)
	}

	client := getClient(config)

	srv, err := tasks.New(client)
	if err != nil {
		return fmt.Errorf("erro ao criar serviço do tasks: %v", err)
	}

	err = srv.Tasks.Delete(taskListID, id).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("erro ao remover tarefa: %v", err)
	}

	return nil
}
