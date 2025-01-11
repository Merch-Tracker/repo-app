package types

type Repo interface {
	Migrate(any) error

	Create(any) error
	CreateOrRewrite(any) error
	Save(any) error

	Read(any, any) error
	ReadOne(any, map[string]any) error
	ReadOnePayload(any, any, map[string]any) error
	ReadManySimple(any, map[string]any) error
	ReadMany(any, map[string]any) error
	ReadManySubmodel(any, any, map[string]any) error

	Update(any, map[string]any) error

	Delete(any, map[string]any) error
}
