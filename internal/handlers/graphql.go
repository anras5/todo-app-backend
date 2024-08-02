package handlers

import (
	"errors"
	"time"

	"github.com/anras5/todo-app-backend/internal/models"
	"github.com/graphql-go/graphql"
)

var TodoType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Todo",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.Int},
			"name":        &graphql.Field{Type: graphql.String},
			"description": &graphql.Field{Type: graphql.String},
			"deadline":    &graphql.Field{Type: graphql.DateTime},
			"completed":   &graphql.Field{Type: graphql.Boolean},
		},
	},
)

type Graph struct {
	QueryString    string
	Variables      map[string]interface{}
	Config         graphql.SchemaConfig
	queryFields    graphql.Fields
	mutationFields graphql.Fields
}

func NewGraph() *Graph {
	var queryFields = graphql.Fields{
		"getTodos": &graphql.Field{
			Type:        graphql.NewList(TodoType),
			Description: "Get all todos",
			Resolve: func(p graphql.ResolveParams) (any, error) {
				todos, err := Repo.DB.SelectTodos()
				if err != nil {
					return nil, err
				}
				return todos, nil
			},
		},
		"getTodo": &graphql.Field{
			Type:        TodoType,
			Description: "Get todo by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					todo, err := Repo.DB.SelectTodo(id)
					if err != nil {
						return nil, err
					}
					return todo, nil
				}
				return nil, errors.New("did not provide id")
			},
		},
	}

	var mutationFields = graphql.Fields{
		"createTodo": &graphql.Field{
			Type:        TodoType,
			Description: "create a new todo",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"deadline": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.DateTime),
				},
				"completed": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Boolean),
				},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				todo := &models.Todo{
					Name:        p.Args["name"].(string),
					Description: p.Args["description"].(string),
					Deadline:    p.Args["deadline"].(time.Time),
					Completed:   p.Args["completed"].(bool),
				}
				id, err := Repo.DB.InsertTodo(*todo)
				if err != nil {
					return nil, err
				}
				todo.ID = id
				return todo, nil
			},
		},
		"updateTodo": &graphql.Field{
			Type:        TodoType,
			Description: "update todo by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"deadline": &graphql.ArgumentConfig{
					Type: graphql.DateTime,
				},
				"completed": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				id, _ := p.Args["id"].(int)
				todo, err := Repo.DB.SelectTodo(id)
				if err != nil {
					return nil, err
				}
				if name, nameOk := p.Args["name"].(string); nameOk {
					todo.Name = name
				}
				if description, descriptionOk := p.Args["description"].(string); descriptionOk {
					todo.Description = description
				}
				if deadline, deadlineOk := p.Args["deadline"].(time.Time); deadlineOk {
					todo.Deadline = deadline
				}
				if completed, completedOk := p.Args["completed"].(bool); completedOk {
					todo.Completed = completed
				}
				err = Repo.DB.UpdateTodo(*todo)
				if err != nil {
					return nil, err
				}
				return todo, nil
			},
		},
		"deleteTodo": &graphql.Field{
			Type:        TodoType,
			Description: "delete todo by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				id, _ := p.Args["id"].(int)
				todo, err := Repo.DB.SelectTodo(id)
				if err != nil {
					return nil, err
				}
				err = Repo.DB.DeleteTodo(id)
				if err != nil {
					return nil, err
				}
				return todo, nil
			},
		},
	}

	return &Graph{
		queryFields:    queryFields,
		mutationFields: mutationFields,
	}
}

func (g *Graph) Query() (*graphql.Result, error) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: g.queryFields,
	})
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: g.mutationFields,
	})

	schemaConfig := graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, err
	}
	params := graphql.Params{
		Schema:         schema,
		RequestString:  g.QueryString,
		VariableValues: g.Variables,
	}
	response := graphql.Do(params)
	if len(response.Errors) > 0 {
		return nil, response.Errors[0]
	}
	return response, nil
}
