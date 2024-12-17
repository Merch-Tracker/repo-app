package types

type Repo interface {
	Migrate(any) error
	Create(any) error
	ReadOne(any, map[string]any) error
	Update(any, map[string]any) error
	Delete(any, map[string]any) error
}
