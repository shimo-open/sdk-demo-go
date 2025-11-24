package consts

type Permission string

const (
	PermissionReadable               Permission = "readable"
	PermissionCopyable               Permission = "copyable"
	PermissionEditable               Permission = "editable"
	PermissionCommentable            Permission = "commentable"
	PermissionLockable               Permission = "lockable"
	PermissionUnlockable             Permission = "unlockable"
	PermissionExportable             Permission = "exportable"
	PermissionFormFillable           Permission = "formFillable"
	PermissionManageable             Permission = "manageable"
	PermissionCopyablePasteClipboard Permission = "copyablePasteClipboard"
	PermissionAttachmentCopyable     Permission = "attachmentCopyable"
	PermissionCutable                Permission = "cutable"
	PermissionAttachmentPreviewable  Permission = "attachmentPreviewable"
	PermissionAttachmentDownloadable Permission = "attachmentDownloadable"
	PermissionImageDownloadable      Permission = "imageDownloadable"
)

func NewFilePermission() []Permission {
	return []Permission{PermissionCopyablePasteClipboard, PermissionAttachmentCopyable, PermissionCutable, PermissionAttachmentPreviewable, PermissionAttachmentDownloadable, PermissionImageDownloadable}
}

func InitFilePermission() []Permission {
	return []Permission{PermissionReadable, PermissionCopyable, PermissionEditable, PermissionCommentable, PermissionLockable, PermissionUnlockable, PermissionExportable, PermissionManageable, PermissionFormFillable}
}

func AllFilePermission() []Permission {
	return []Permission{PermissionReadable, PermissionCopyable, PermissionEditable, PermissionCommentable, PermissionLockable, PermissionUnlockable, PermissionExportable, PermissionManageable, PermissionFormFillable, PermissionCopyablePasteClipboard, PermissionAttachmentCopyable, PermissionCutable, PermissionAttachmentPreviewable, PermissionAttachmentDownloadable, PermissionImageDownloadable}
}

func HandleFilePermission(b bool) map[string]bool {
	// formFillable is handled elsewhere
	return map[string]bool{
		string(PermissionReadable):               b,
		string(PermissionCopyable):               b,
		string(PermissionEditable):               b,
		string(PermissionCommentable):            b,
		string(PermissionLockable):               b,
		string(PermissionUnlockable):             b,
		string(PermissionExportable):             b,
		string(PermissionManageable):             b,
		string(PermissionCopyablePasteClipboard): b,
		string(PermissionAttachmentCopyable):     b,
		string(PermissionCutable):                b,
		string(PermissionAttachmentPreviewable):  b,
		string(PermissionAttachmentDownloadable): b,
		string(PermissionImageDownloadable):      b,
	}
}

func HandleBasicFilePermission(b bool) map[string]bool {
	// Populate the core permission set
	return map[string]bool{
		string(PermissionReadable):    b,
		string(PermissionCopyable):    b,
		string(PermissionEditable):    b,
		string(PermissionCommentable): b,
		string(PermissionLockable):    b,
		string(PermissionUnlockable):  b,
		string(PermissionExportable):  b,
		string(PermissionManageable):  b,
	}
}
