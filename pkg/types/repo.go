package types

type Repo interface {
	Migrate(any) error
	Create(any) error
	ReadOne(any, map[string]any) error
	ReadOnePayload(any, any, map[string]any) error
	ReadMany(any, map[string]any) error
	Update(any, map[string]any) error
	Delete(any, map[string]any) error
}
