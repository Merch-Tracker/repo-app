package merch

const (
	errMsg   = "error"
	bytesMsg = "bytes sent"

	merchTableError      = "Merchandise table migration failed"
	pricesTableError     = "Prices table migration failed"
	labelsTableError     = "Labels table migration failed"
	cardLabelsTableError = "Card labels table migration failed"
	migrationsSuccess    = "All migrations are successful"

	newMerchValidationError = "Validation failed"

	merchCreateSuccess = "Merchandise create successfully"
	merchCreateError   = "Merchandise creation failed"
	merchReadSuccess   = "Merchandise read successfully"
	merchReadError     = "Merchandise read failed"
	merchUpdateSuccess = "Merchandise update successfully"
	merchUpdateError   = "Merchandise update failed"
	merchDeleteSuccess = "Merchandise delete successfully"
	merchDeleteError   = "Merchandise delete failed"

	labelsCreateError   = "Create new label failed"
	labelsCreateSuccess = "New label created"
	labelsGetAllError   = "Get all labels failed"
	labelsGetSuccess    = "Get label success"
	labelsGetIdError    = "Get label id failed"
	labelsUpdateError   = "Update label failed"
	labelsUpdateSuccess = "Update label successfully"
	labelsDeleteError   = "Delete label failed"
	labelsDeleteSuccess = "Delete label successfully"

	labelAttachError   = "Attach label failed"
	labelAttachSuccess = "Attach label successfully"
	labelDetachError   = "Detach label failed"
	labelDetachSuccess = "Detach label successfully"

	urlParseError            = "URL parse failed"
	queryParseError          = "Query params parse failed"
	getPriceHistoryError     = "Price history read failed"
	serPriceHistory          = "Serializing price history failed"
	priceHistoryFetchSuccess = "Price history fetched successfully"
)
