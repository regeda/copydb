package copydb

import "fmt"

const (
	keyVer = "$ver"
	keyID  = "$id"
)

var defaultKeys = keys{
	itemPattern: "{%s}:item",
	list:        "items_list",
	channel:     "items_update",
}

type keys struct {
	itemPattern, list, channel string
}

func (k *keys) item(id string) string {
	return fmt.Sprintf(k.itemPattern, id)
}
