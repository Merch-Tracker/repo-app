package types

type Repo interface {
	Migrate() error
}
