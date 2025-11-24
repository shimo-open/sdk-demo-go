export const fileConstants = {
  CLEAR_ERROR: 'CLEAR_ERROR',

  GET_FILES_REQUEST: 'GET_FILES_REQUEST',
  GET_FILES_SUCCESS: 'GET_FILES_SUCCESS',
  GET_FILES_FAILURE: 'GET_FILES_FAILURE',

  CREATE_SHIMO_FILE_REQUEST: 'CREATE_SHIMO_FILE_REQUEST',
  CREATE_SHIMO_FILE_SUCCESS: 'CREATE_SHIMO_FILE_SUCCESS',
  CREATE_SHIMO_FILE_FAILURE: 'CREATE_SHIMO_FILE_FAILURE',
  CREATE_SHIMO_FILE_RESET: 'CREATE_SHIMO_FILE_RESET',

  UPLOAD_FILE_REQUEST: 'UPLOAD_FILE_REQUEST',
  UPLOAD_FILE_SUCCESS: 'UPLOAD_FILE_SUCCESS',
  UPLOAD_FILE_FAILURE: 'UPLOAD_FILE_FAILURE',
  UPLOAD_FILE_RESET: 'UPLOAD_FILE_RESET',

  GET_FILE_REQUEST: 'GET_FILE_REQUEST',
  GET_FILE_SUCCESS: 'GET_FILE_SUCCESS',
  GET_FILE_FAILURE: 'GET_FILE_FAILURE',
  GET_FILE_RESET: 'GET_FILE_RESET',

  CLEAR_FILE_CACHE: 'CLEAR_FILE_CACHE',
  CLEAR_IMPORT_CACHE: 'CLEAR_IMPORT_CACHE',

  IMPORT_FILE_REQUEST: 'IMPORT_FILE_REQUEST',
  IMPORT_FILE_SUCCESS: 'IMPORT_FILE_SUCCESS',
  IMPORT_FILE_FAILURE: 'IMPORT_FILE_FAILURE',
  IMPORT_FILE_PROGRESS: 'IMPORT_FILE_PROGRESS',

  IMPORT_URL_REQUEST: 'IMPORT_URL_REQUEST',
  IMPORT_URL_SUCCESS: 'IMPORT_URL_SUCCESS',
  IMPORT_URL_FAILURE: 'IMPORT_URL_FAILURE',

  EXPORT_FILE_REQUEST: 'EXPORT_FILE_REQUEST',
  EXPORT_FILE_SUCCESS: 'EXPORT_FILE_SUCCESS',
  EXPORT_FILE_FAILURE: 'EXPORT_FILE_FAILURE',
  EXPORT_FILE_PROGRESS: 'EXPORT_FILE_PROGRESS',
  EXPORT_FILE_CLEAR: 'EXPORT_FILE_CLEAR',
  EXPORT_FILE_URL_PREFIX: 'EXPORT_FILE_URL_PREFIX',
  EXPORT_FILE_PROGRESS_PREFIX: 'EXPORT_FILE_PROGRESS_PREFIX',

  REMOVE_FILE_REQUEST: 'REMOVE_FILE_REQUEST',
  REMOVE_FILE_SUCCESS: 'REMOVE_FILE_SUCCESS',
  REMOVE_FILE_FAILURE: 'REMOVE_FILE_FAILURE',

  DUPLICATE_FILE_REQUEST: 'DUPLICATE_FILE_REQUEST',
  DUPLICATE_FILE_SUCCESS: 'DUPLICATE_FILE_SUCCESS',
  DUPLICATE_FILE_FAILURE: 'DUPLICATE_FILE_FAILURE',

  TYPE_DOCUMENT: 'document',
  TYPE_DOCUMENT_PRO: 'documentPro',
  TYPE_SPREADSHEET: 'spreadsheet',
  TYPE_PRESENTATION: 'presentation',
  TYPE_TABLE: 'table',
  TYPE_FORM: 'form',
  TYPE_BOARD: 'board',
  TYPE_MINDMAP: 'mindmap',
  TYPE_FLOWCHART: 'flowchart',

  GET_ALL_COLLABORATORS_REQUEST: 'GET_ALL_COLLABORATORS_REQUEST',
  GET_ALL_COLLABORATORS_SUCCESS: 'GET_ALL_COLLABORATORS_SUCCESS',
  GET_ALL_COLLABORATORS_FAILURE: 'GET_ALL_COLLABORATORS_FAILURE',

  GET_COMMENT_COUNT_REQUEST: 'GET_COMMENT_COUNT_REQUEST',
  GET_COMMENT_COUNT_SUCCESS: 'GET_COMMENT_COUNT_SUCCESS',
  GET_COMMENT_COUNT_FAILURE: 'GET_COMMENT_COUNT_FAILURE',

  GET_MENTION_AT_REQUEST: 'GET_MENTION_AT_REQUEST',
  GET_MENTION_AT_SUCCESS: 'GET_MENTION_AT_SUCCESS',
  GET_MENTION_AT_FAILURE: 'GET_MENTION_AT_FAILURE',

  GET_REVISIONS_REQUEST: 'GET_REVISIONS_REQUEST',
  GET_REVISIONS_SUCCESS: 'GET_REVISIONS_SUCCESS',
  GET_REVISIONS_FAILURE: 'GET_REVISIONS_FAILURE',

  REVISIONS_MODAL_OPEN_STATUS_CHANGED: 'REVISIONS_MODAL_OPEN_STATUS_CHANGED',

  GET_HISTORIES_REQUEST: 'GET_HISTORIES_REQUEST',
  GET_HISTORIES_SUCCESS: 'GET_HISTORIES_SUCCESS',
  GET_HISTORIES_FAILURE: 'GET_HISTORIES_FAILURE',

  HISTORIES_MODAL_OPEN_STATUS_CHANGED: 'HISTORIES_MODAL_OPEN_STATUS_CHANGED',

  UPDATE_FILE: 'UPDATE_FILE',

  SAVE_STATUS_CHANGED: 'SAVE_STATUS_CHANGED',
  SAVE_STATUS_INIT: 'SAVE_STATUS_INIT',

  TYPES: {
    'x-compressed': '压缩文件',
    'x-zip-compressed': '压缩文件',
    zip: '压缩文件',
    'x-zip': '压缩文件',
    'x-7z-compressed': '压缩文件',
    'x-rar-compressed': '压缩文件',
    msword: 'Word 文件',
    'vnd.openxmlformats-officedocument.wordprocessingml.document': 'Word 文件',
    'vnd.oasis.opendocument.text': 'Word 文件',
    'vnd.ms-powerpoint': 'Powerpoint 文件',
    'vnd.openxmlformats-officedocument.presentationml.presentation':
      'Powerpoint 文件',
    'vnd.ms-excel': 'Excel 文件',
    'vnd.openxmlformats-officedocument.spreadsheetml.sheet': 'Excel 文件',
    pdf: 'PDF',
    kswps: 'WPS 文件',
    'text/plain': '纯文本文件',
    'text/markdown': 'Markdown 文件',
    'application/rtf': '多信息文本文件',
    'text/rtf': '多信息文本文件',
    'text/csv': 'CSV 文件',
    'vnd.xmind': 'Xmind 文件'
  },

  PREVIEWABLE_EXTNAMES: [
    // Document formats
    '.doc', '.docx', '.wps', '.wpt', '.xls', '.xlsx', '.csv', '.xlsm',
    '.ppt', '.pptx', '.xmind', '.pdf', '.ofd', '.rtf',
    '.txt', '.markdown', '.md', '.log', '.ini', '.conf', '.adoc',

    // Image formats
    '.jpeg', '.jpg', '.png', '.gif', '.bmp', '.svg', '.webp',
    '.heic', '.heif', '.avif',

    // Video formats
    '.mp4', '.avi', '.flv', '.mpeg', '.webm', '.mov',

    // Audio formats
    '.mp3', '.ogg', '.flac', '.aac', '.m4a', '.wav', '.oga',

    // Code / configuration formats
    '.1c', '.abnf', '.as', '.ada', '.angelscript', '.applescript', '.arcade',
    '.ino', '.s', '.xml', '.aj', '.ahk', '.au3', '.asm', '.awk', '.axapta',
    '.sh', '.bas', '.bnf', '.bf', '.c', '.cal', '.capnp', '.ceylon', '.icl',
    '.clj', '.cmake', '.coffee', '.coq', '.cos', '.cpp', '.crmsh', '.cr',
    '.cs', '.csp', '.css', '.d', '.dart', '.pas', '.diff', '.py', '.zone',
    '.dockerfile', '.bat', '.ldif', '.d.ts', '.dust', '.ebnf', '.ex', '.elm',
    '.rb', '.erb', '.erl', '.fix', '.flix', '.f', '.fs', '.gms', '.gss',
    '.gcode', '.feature', '.glsl', '.gml', '.go', '.golo', '.gradle',
    '.graphql', '.groovy', '.haml', '.hbs', '.hs', '.hx', '.hsp', '.http',
    '.hy', '.inform7', '.irpf90', '.isbl', '.java', '.js', '.cli', '.json',
    '.jl', '.kt', '.lasso', '.tex', '.leaf', '.less', '.lisp',
    '.livecodeserver', '.ls', '.ll', '.lsl', '.lua', '.makefile', '.m', '.max',
    '.mel', '.mizar', '.pl', '.monkey', '.moon', '.n1ql', '.nt', '.nim',
    '.nix', '.nsi', '.ml', '.scad', '.oxygene', '.parser3', '.pf', '.sql',
    '.php', '.pony', '.ps1', '.pde', '.profile', '.properties', '.proto',
    '.pp', '.pb', '.q', '.qml', '.r', '.re', '.rib', '.roboconf', '.rsc',
    '.rsl', '.ruleslanguage', '.rs', '.sas', '.scala', '.scm', '.sci',
    '.scss', '.smali', '.st', '.sml', '.sqf', '.stan', '.do', '.step21',
    '.styl', '.subunit', '.swift', '.taggerscript', '.yml', '.tap', '.tcl',
    '.thrift', '.tp', '.twig', '.ts', '.vala', '.vb', '.vbs', '.v', '.vhd',
    '.vim', '.wasm', '.wren', '.xl', '.xq', '.zep'
  ],

  PREVIEWABLE_MIME_TYPES: [
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/vnd.ms-excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'application/vnd.ms-powerpoint',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation',
    'application/pdf',
    'text/plain',
    'image/jpeg',
    'image/jpeg',
    'image/png',
    'image/gif',
    'image/bmp',
    'image/svg+xml',
    'video/mp4',
    'audio/mpeg',
    'text/markdown',
    'application/rtf',
    'text/rtf',
    'text/csv',
    'application/vnd.ms-works', // WPS
    'x-lml/x-gps', // WPT
    'application/wpsoffice', // WPS,WPT
    'image/heic',
    'image/heif',
    'application/vnd.xmind'
  ],

  PERMISSIONS: [
    { label: '可查看', value: 'readable' },
    { label: '可复制', value: 'copyable' },
    { label: '可评论', value: 'commentable' },
    { label: '可编辑', value: 'editable' },
    { label: '可导出', value: 'exportable' },
    { label: '可管理', value: 'manageable' },
  ],

  NEW_PERMISSIONS: [
    { label: '可查看', value: 'readable' },
    { label: '可复制', value: 'copyable' },
    { label: '可外部粘贴', value: 'copyablePasteClipboard' },
    { label: '可评论', value: 'commentable' },
    { label: '可编辑', value: 'editable' },
    { label: '可导出', value: 'exportable' },
    { label: '可剪切', value: 'cutable' },
    { label: '可复制附件', value: 'attachmentCopyable' },
    { label: '可预览附件', value: 'attachmentPreviewable' },
    { label: '可下载附件', value: 'attachmentDownloadable' },
    { label: '可下载图片', value: 'imageDownloadable' },
    { label: '可管理', value: 'manageable' },
  ]
}
