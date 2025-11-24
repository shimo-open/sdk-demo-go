package consts

type FileType string

const (
	// FileTypeDocument represents a lite document
	FileTypeDocument FileType = "document"
	// FileTypeSpreadsheet represents a spreadsheet (Excel)
	FileTypeSpreadsheet FileType = "spreadsheet"
	// FileTypeDocPro represents a document pro (Word)
	FileTypeDocPro FileType = "documentPro"
	// FileTypeSlide represents a slide deck (PPT)
	FileTypeSlide FileType = "presentation"
	// FileTypeTable represents an application table
	FileTypeTable FileType = "table"
	// FileTypeInvalid represents an invalid type
	FileTypeInvalid FileType = "invalid"
	// FileTypeAll represents all types
	FileTypeAll FileType = "all"
)

func (f FileType) String() string {
	return string(f)
}

func GetFileType(ft string) FileType {
	for _, t := range AllFileTypes() {
		if FileType(ft) == t {
			return t
		}
	}

	return FileTypeInvalid
}

func AllFileTypes() []FileType {
	return []FileType{FileTypeDocument, FileTypeSpreadsheet, FileTypeDocPro, FileTypeSlide, FileTypeTable}
}

var (
	ImportTypeMap = map[FileType][]string{
		FileTypeDocument:    {"docx", "doc", "md", "txt"},
		FileTypeDocPro:      {"docx", "doc", "wps"},
		FileTypeSpreadsheet: {"xlsx", "xls", "csv", "xlsm"},
		FileTypeSlide:       {"pptx", "ppt"},
	}
	ExportTypeMap = map[FileType][]string{
		FileTypeDocument:    {"docx", "md", "jpg", "pdf"},
		FileTypeDocPro:      {"docx", "pdf", "wps"},
		FileTypeSpreadsheet: {"xlsx"},
		FileTypeSlide:       {"pptx", "pdf"},
	}
)
