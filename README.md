✨ bcli - Gerenciador de Tarefas no Terminal

O bcli é um gerenciador de tarefas interativo em modo texto (TUI), feito em Go com Bubble Tea
 e Lipgloss
.
Ele permite gerenciar tarefas localmente e, opcionalmente, sincronizá-las com o Google Tasks.

🚀 Funcionalidades

✅ Adicionar, listar, editar, concluir e remover tarefas

🎨 Interface interativa e colorida com Bubble Tea + Lipgloss

☁️ Integração opcional com Google Tasks

💾 Persistência simples em memória (local) + suporte a salvar/editar

🖥️ Totalmente em linha de comando

📦 Instalação
Pré-requisitos

Go 1.22+

Conta Google (opcional, apenas se quiser sincronizar com Google Tasks)

Clonar o repositório
git clone https://github.com/seunome/bcli.git
cd bcli

Rodar o programa
go run main.go

Compilar para binário
go build -o bcli main.go
./bcli

📝 Como usar

Ao iniciar o bcli, você verá uma tela de boas-vindas.
Pressione qualquer tecla (ou aguarde 3 segundos) para começar.

🔑 Comandos disponíveis
Comando	Descrição	Exemplo
help	Mostra todos os comandos disponíveis	help
list	Lista todas as tarefas	list
add <título>	Adiciona uma nova tarefa	add Estudar Go
done <id>	Marca/desmarca uma tarefa como concluída	done 2
remove <id>	Remove uma tarefa pelo ID	remove 3
edit <id> <novo>	Edita o título de uma tarefa	edit 1 Revisar documentação
quit / q / esc	Sai do programa	quit
📊 Exemplo de uso
✨ bcli - Todo List Interativo ✨

bcli> add Comprar pão
Tarefa adicionada: [ ] 1 Comprar pão

bcli> add Estudar Go
Tarefa adicionada: [ ] 2 Estudar Go

bcli> list
[ ] 1 Comprar pão
[ ] 2 Estudar Go

bcli> done 1
Tarefa 1 marcada como feita/não feita.

bcli> list
[x] 1 Comprar pão
[ ] 2 Estudar Go

☁️ Integração com Google Tasks (Opcional)

O bcli pode sincronizar suas tarefas com o Google Tasks.
Para isso, configure as credenciais da API do Google:

Acesse Google Cloud Console

Crie um projeto e habilite a API Google Tasks

Gere as credenciais OAuth2 (JSON)

Coloque o arquivo de credenciais na pasta ~/.bcli/credentials.json

Rode o programa normalmente (go run main.go) e siga o fluxo de autenticação.

⚠️ Essa parte ainda pode variar dependendo da sua implementação em internal/google-tasks.

📂 Estrutura do projeto
bcli/
├── main.go                 # Programa principal (interface interativa)
├── internal/
│   ├── tasks/              # Serviço local de tarefas (TaskService)
│   │   └── tasks.go
│   └── google-tasks/       # Integração com Google Tasks
│       └── ...
└── go.mod

🛠 Tecnologias usadas

Go
 – Linguagem principal

Bubble Tea
 – Framework TUI

Lipgloss
 – Estilização

Google Tasks API
 (opcional)

📜 Licença

Este projeto é distribuído sob a licença MIT.
Sinta-se livre para usar, modificar e contribuir.

🤝 Contribuindo

Pull requests são bem-vindos!
Se tiver ideias de melhorias, abra uma issue ou mande seu PR.