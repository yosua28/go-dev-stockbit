package models

type MsFile struct {
	FileKey            uint64             `db:"file_key"                json:"file_key"`
	RefFkKey           uint64             `db:"ref_fk_key"              json:"ref_fk_key"`
	RefFkDomain        string             `db:"ref_fk_domain"           json:"ref_fk_domain"`
	FileName           string             `db:"file_name"               json:"file_name"`
	FileExt            string             `db:"file_ext"                json:"file_ext"`
	BlobMode           uint8              `db:"blob_mode"               json:"blob_mode"`
	FilePath           *string            `db:"file_path"               json:"file_path"`
	FileUrl            *string            `db:"file_url"                json:"file_url"`
	FileNotes          *string            `db:"file_notes"              json:"file_notes"`
	FileObj            *interface{}       `db:"file_obj"                json:"properties"`
	RecCreatedDate     *string            `db:"rec_created_date"        json:"rec_created_date"`
	RecCreatedBy       *string            `db:"rec_created_by"          json:"rec_created_by"`
	RecModifiedDate    *string            `db:"rec_modified_date"       json:"rec_modified_date"`
	RecModifiedBy      *string            `db:"rec_modified_by"         json:"rec_modified_by"`
}