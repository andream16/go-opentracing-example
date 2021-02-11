package todo

// Todo describes a todo.
type Todo struct {
	Message string `json:"message"`
}

// New returns a new Todo.
func New(message string) Todo {
	return Todo{Message: message}
}
