package doc

import (
	"context"

	api "diploma-project/api/web/gen"
)

func (Handler) Doc(
	_ context.Context,
	_ api.DocRequestObject,
) (api.DocResponseObject, error) {
	spec, _ := api.GetSwagger()
	return api.Doc200ApplicationJSONCharsetUTF8Response(*spec), nil
}
