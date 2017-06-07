package model

// DBStats ...
type DBStats struct {
	DB          string `bson:"db" json:"db"`
	AvgObjSize  *int   `bson:"avgObjSize,omitempty" json:"avgObjSize,omitempty"`
	Collections *int   `bson:"collections,omitempty" json:"collections,omitempty"`
	DataSize    *int   `bson:"dataSize,omitempty" json:"dataSize,omitempty"`
	FileSize    *int   `bson:"fileSize,omitempty" json:"fileSize,omitempty"`
	IndexSize   *int   `bson:"indexSize,omitempty" json:"indexSize,omitempty"`
	Indexes     *int   `bson:"indexes,omitempty" json:"indexes,omitempty"`
	NumExtents  *int   `bson:"numExtents,omitempty" json:"numExtents,omitempty"`
	Objects     *int   `bson:"objects,omitempty" json:"objects,omitempty"`
	StorageSize *int   `bson:"storageSize,omitempty" json:"storageSize,omitempty"`
	OK          int    `bson:"ok" json:"ok"`
	Error       string `json:"error,omitempty"`
}
