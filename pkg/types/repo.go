package types

type Repo interface {
	Migrate(any) error
	Create(any) error
}
