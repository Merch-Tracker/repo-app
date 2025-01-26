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

	ReadCharts(any, map[string]any) error
	ReadRaw(string, any) error

	Update(any, map[string]any) error
	UpdateNotifications(any, []uint, string) error

	Delete(any, map[string]any) error
}
