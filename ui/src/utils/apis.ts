export const apis = {
    "apis": [
        {
            "sdk": "all",
            "methods": [
                {
                    "method": "setTitle",
                    "description": "设置文档标题",
                    "support mobile": "Y",
                    "params": [
                        {
                            "param name": "title",
                            "required": "Y",
                            "param type": "string",
                            "param description": "用于实现客户侧更新标题后，同步到编辑器中，比如打印时需要显示最新的标题",
                        }
                    ],
                }
            ]
        },
        {
            "sdk": "document",
            "methods": [
                {
                    "method": "showHistory",
                    "description": "显示历史侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "hideHistory",
                    "description": "隐藏历史侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "showRevision",
                    "description": "显示版本侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "hideRevision",
                    "description": "隐藏版本侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "showDiscussion",
                    "description": "显示讨论侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "hideDiscussion",
                    "description": "隐藏讨论侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "showComments",
                    "description": "显示评论侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "hideComments",
                    "description": "隐藏评论侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "showToc",
                    "description": "显示目录",
                    "support mobile": "N",
                },
                {
                    "method": "hideToc",
                    "description": "隐藏目录",
                    "support mobile": "N",
                },
                {
                    "method": "createRevision",
                    "description": "创建版本",
                    "support mobile": "N",
                    "returns": [
                        {
                            "in/out": "out",
                            "param name": "errorMessage",
                            "param type": "{ message: string }",
                            "param description": "创建版本出错时reject，错误信息"
                        }
                    ],
                },
                {
                    "method": "startDemonstration",
                    "description": "进入演示模式",
                    "support mobile": "N",
                },
                {
                    "method": "endDemonstration",
                    "description": "退出演示模式",
                    "support mobile": "N",
                },
                {
                    "method": "print",
                    "description": "打印",
                    "support mobile": "N",
                },
                {
                    "method": "getTitle",
                    "description": "获取文档标题",
                    "support mobile": "Y",
                },
            ]
        },
        {
            "sdk": "documentPro",
            "methods": [
                {
                    "method": "getComments",
                    "description": "获取所有评论",
                    "support mobile": "Y",
                    "params": [
                        {
                            "param name": "options",
                            "param type": "{\n  /* 包含对应的标题信息 */\n  includeChapterTitle: boolean\n}",
                            "default": "false",
                        }
                    ],
                    "returns": [
                        {
                            "in/out": "out",
                            "param name": "comments",
                            "param type": "documentPro.Comment[]",
                            "param description": "文档中的所有评论",
                        }
                    ],
                },
                {
                    "method": "getComment",
                    "description": "获取单条评论",
                    "support mobile": "Y",
                    "params": [
                        {
                            "param name": "commentId",
                            "required": "Y",
                            "param type": "string",
                            "param description": "评论ID",
                        },
                        {
                            "param name": "includeChapterTitle",
                            "param type": "boolean",
                            "param description": "包含对应的标题信息"
                        }
                    ],
                    "returns": [
                        {
                            "in/out": "out",
                            "param name": "comment",
                            "param type": "documentPro.Comment",
                            "param description": "评论"
                        }
                    ],
                },
                {
                    "method": "showToc",
                    "description": "显示文档结构",
                    "support mobile": "Y",
                    "params": [
                        {
                            "param name": "shouldDisableCache",
                            "param type": "boolean",
                            "param description": "禁用文档结构图的默认缓存",
                            "default": "true",
                        },
                        {
                            "param name": "collapsedLevel",
                            "param type": "number",
                            "param description": "默认折叠的层级",
                        },
                        {
                            "param name": "itemHeight",
                            "param type": "number",
                            "param description": "文档结构图的行高, 例如 24",
                        }
                    ],
                },
                {
                    "method": "hideToc",
                    "description": "隐藏文档结构图",
                    "support mobile": "Y",
                },
                {
                    "method": "print",
                    "description": "打印",
                    "support mobile": "Y",
                    "params": [
                        {
                            "param name": "pageNums",
                            "required": "Y",
                            "param type": "number[]",
                            "param description": "页码列表",
                        }
                    ],
                },
                {
                    "method": "printAll",
                    "description": "打印所有页面",
                    "support mobile": "Y",
                },
                {
                    "method": "createRevision",
                    "description": "创建版本",
                    "support mobile": "Y",
                    "params": [
                        {
                            "param name": "name",
                            "required": "Y",
                            "param type": "string",
                            "param description": "版本名",
                        }
                    ],
                },
                {
                    "method": "showHistory",
                    "description": "预览历史版本",
                    "support mobile": "Y",
                },
                {
                    "method": "hideHistory",
                    "description": "关闭历史版本预览",
                    "support mobile": "Y",
                },
            ]
        },
        {
            "sdk": "sheet",
            "methods": [
                {
                    "method": "showComments",
                    "description": "展示评论侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "hideComments",
                    "description": "关闭评论侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "showHistory",
                    "description": "展示历史侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "hideHistory",
                    "description": "关闭历史侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "hideLocks",
                    "description": "关闭锁定侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "createRevision",
                    "description": "创建版本",
                    "support mobile": "N",
                },
                {
                    "method": "startDemonstration",
                    "description": "进入演示模式",
                    "support mobile": "N",
                },
                {
                    "method": "endDemonstration",
                    "description": "离开演示模式",
                    "support mobile": "N",
                },
                {
                    "method": "print",
                    "description": "打印",
                    "support mobile": "N",
                },
            ]
        },
        {
            "sdk": "presentation",
            "methods": [
                {
                    "method": "showHistory",
                    "description": "显示历史",
                    "support mobile": "N",
                },
                {
                    "method": "hideHistory",
                    "description": "隐藏历史",
                    "support mobile": "N",
                },
                {
                    "method": "startDemonstration",
                    "description": "开始本地演示",
                    "support mobile": "N",
                },
                {
                    "method": "endDemonstration",
                    "description": "结束本地演示",
                    "support mobile": "N",
                },
                {
                    "method": "print",
                    "description": "打印",
                    "support mobile": "N",
                },
                {
                    "method": "setContent",
                    "description": "设置文件内容",
                    "support mobile": "N",
                    "params": [
                        {
                            "param name": "content",
                            "required": "Y",
                            "param type": "any",
                            "param description": "要设置的文件内容，会替换当前内容，实际类型接受 string | Delta",
                        }
                    ],
                },
            ]
        },
        {
            "sdk": "table",
            "methods": [
                {
                    "method": "showRevision",
                    "description": "显示版本侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "hideRevision",
                    "description": "隐藏版本侧边栏",
                    "support mobile": "N",
                },
                {
                    "method": "createRevision",
                    "description": "创建版本",
                    "support mobile": "N",
                },
            ]
        }
    ]
}