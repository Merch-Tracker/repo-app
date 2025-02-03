package types

type Repo interface {
	Migrate(any) error

	Create(any) error
	CreateOrRewrite(any) error
	CreateWithTransaction(any, any, any) error

	Read(any, any) error
	ReadOne(any, map[string]any) error
	ReadOnePayload(any, any, map[string]any) error
	ReadManySimple(any, map[string]any) error
	ReadManySimpleSubmodel(any, any, map[string]any) error
	ReadManyInList(any, []string) error
	ReadPrices(any, []string, int) error
	ReadManySubmodel(any, any, map[string]any) error

	ReadCharts(any, map[string]any) error
	ReadRaw(string, any) error

	Update(any, map[string]any) error
	UpdateWithTransaction(any, any, map[string]any) error
	UpdateNotifications(any, []uint, string) error

	Delete(any, map[string]any) error
	DeleteWithTransaction(any, any, map[string]any) error
}
