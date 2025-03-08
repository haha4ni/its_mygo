package database

type ImageEntry struct {
	Filename string
	ImgIndex int
	Page     int
}

// 實作 Storable 介面
func (img ImageEntry) TableName() string {
	return "images"
}

func (img ImageEntry) Fields() map[string]any {
	return map[string]any{
		"filename":  img.Filename,
		"img_index": img.ImgIndex,
		"page":      img.Page,
	}
}