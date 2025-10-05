package common

type RepositoryType string

const (
	TypeMemory   RepositoryType = "memory"
	TypePostgres RepositoryType = "postgres"
)
