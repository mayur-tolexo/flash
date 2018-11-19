package flash

import "reflect"

var empty = ""

//MetaData contains the api metadata
type MetaData struct {
	prefix  string
	root    string
	url     string
	version string
}

//getMetaData will return metadata from service root
//i.e. the anonymous Server field of the struct
func getMetaData(service interface{}, rootName string) (meta MetaData) {
	meta = MetaData{}
	if rootField, exists := reflect.TypeOf(service).Elem().FieldByName(rootName); exists {
		tag := rootField.Tag
		meta.prefix = coalesce(tag, "prefix", "pre")
		meta.root = coalesce(tag, "root")
		meta.url = coalesce(tag, "url")
		meta.version = coalesce(tag, "version", "v")
	}
	return
}

//getInOrderMetaData will return metadata based on inorder input
func getInOrderMetaData(m ...MetaData) MetaData {
	meta := MetaData{}
	for _, curMeta := range m {
		if meta.prefix == empty && curMeta.prefix != empty {
			meta.prefix = curMeta.prefix
		}
		if meta.root == empty && curMeta.root != empty {
			meta.root = curMeta.root
		}
		if meta.url == empty && curMeta.url != empty {
			meta.url = curMeta.url
		}
		if meta.version == empty && curMeta.version != empty {
			meta.version = curMeta.version
		}
	}
	return meta
}
