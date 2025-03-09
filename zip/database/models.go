package database

type ImageData struct {
	Filename string `db:"filename"`
	ImgIndex int    `db:"img_index"`
	Page     int    `db:"page"`
	ZipID    int    `db:"zip_id"` // 外鍵
}

type ZipFile struct {
	FileName  string      `db:"filename"`
	SHA       string      `db:"sha"`
	Timestamp int64       `db:"timestamp"`
	ImageData []ImageData // 這部分不直接存進 SQL，查詢時額外 JOIN
}

// 實作 Storable 介面
func (zip ZipFile) TableName() string {
	return "zip_files"
}

func (zip ZipFile) Fields() map[string]any {
	return map[string]any{
		"filename":  zip.FileName,
		"sha":       zip.SHA,
		"timestamp": zip.Timestamp,
	}
}

// 實作 Storable 介面
func (img ImageData) TableName() string {
	return "image_data"
}

func (img ImageData) Fields() map[string]any {
	return map[string]any{
		"filename":  img.Filename,
		"img_index": img.ImgIndex,
		"page":      img.Page,
		"zip_id":    img.ZipID,
	}
}