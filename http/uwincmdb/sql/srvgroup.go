package sql

import "fmt"

var QueryName Query = func(name string) string {
	return fmt.Sprintf("\"query\":{\"$and\":[{\"groupId\":{\"$exists\":true}},{\"srvId\":{\"$exists\":true}},{\"name\":{\"$eq\":\"%s\"}}]}", name)
}

var QueryID Query = func(groupID string) string {
	return fmt.Sprintf("\"query\":{\"$and\":[{\"groupId\":{\"$exists\":true}},{\"srvId\":{\"$exists\":true}},{\"groupId\":{\"$eq\":\"%s\"}}]}", groupID)
}
