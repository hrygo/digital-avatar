package ai

// initializeIntentPatterns 初始化意图模式
func initializeIntentPatterns() map[IntentType][]*IntentPattern {
	patterns := make(map[IntentType][]*IntentPattern)

	// 任务意图模式
	patterns[IntentTask] = []*IntentPattern{
		{
			Pattern:   NewRegexpWrapper(`(需要|要|请|麻烦|帮我|协助|支持|安排|处理|做|完成|执行|搞|弄|弄好|搞定|解决)`),
			Weight:    3.0,
			Keywords:  []string{"需要", "要", "请", "帮忙", "协助"},
			Context:   "任务执行请求",
		},
	}

	// 问题意图模式
	patterns[IntentQuestion] = []*IntentPattern{
		{
			Pattern:   NewRegexpWrapper(`(什么|如何|怎么|为什么|是否|能不能|可以吗|有没有|谁知道|请问|麻烦问)`),
			Weight:    2.5,
			Keywords:  []string{"什么", "如何", "怎么", "为什么"},
			Context:   "信息询问",
		},
	}

	// 决策意图模式
	patterns[IntentDecision] = []*IntentPattern{
		{
			Pattern:   NewRegexpWrapper(`(决定|选择|判断|考虑|想|觉得|认为|看法|意见|建议|要不要|应该)`),
			Weight:    2.0,
			Keywords:  []string{"决定", "选择", "判断"},
			Context:   "决策需求",
		},
	}

	// 社交意图模式
	patterns[IntentSocial] = []*IntentPattern{
		{
			Pattern:   NewRegexpWrapper(`(聊聊|谈谈|说说|讨论|交流|分享|约|聚会|见面|一起|一起)`),
			Weight:    1.5,
			Keywords:  []string{"聊聊", "谈谈", "讨论", "交流"},
			Context:   "社交互动",
		},
	}

	// 情感表达模式
	patterns[IntentEmotional] = []*IntentPattern{
		{
			Pattern:   NewRegexpWrapper(`(我感觉|我觉得|心情|情绪|难过|开心|高兴|生气|担心|害怕|紧张|焦虑)`),
			Weight:    1.8,
			Keywords:  []string{"感觉", "心情", "情绪"},
			Context:   "情感表达",
		},
	}

	// 请求意图模式
	patterns[IntentRequest] = []*IntentPattern{
		{
			Pattern:   NewRegexpWrapper(`(可以|能|能不能|麻烦你|请|求|希望|想要|需要)`),
			Weight:    2.5,
			Keywords:  []string{"可以", "能不能", "麻烦你"},
			Context:   "请求协助",
		},
	}

	return patterns
}

// initializeIntentKeywords 初始化意图关键词
func initializeIntentKeywords() map[IntentType][]string {
	return map[IntentType][]string{
		IntentTask: {"需要", "要", "请", "麻烦", "帮我", "协助", "支持", "安排", "处理", "做", "完成", "执行", "搞", "弄", "弄好", "搞定", "解决", "建立", "创建"},
		IntentQuestion: {"什么", "如何", "怎么", "为什么", "是否", "能不能", "可以吗", "有没有", "谁知道", "请问", "麻烦问", "想问"},
		IntentDecision: {"决定", "选择", "判断", "考虑", "想", "觉得", "认为", "看法", "意见", "建议", "要不要", "应该", "还是", "或者是"},
		IntentInfo: {"告诉", "说明", "介绍", "解释", "展示", "提供", "分享", "传达", "通知"},
		IntentSocial: {"聊聊", "谈谈", "说说", "讨论", "交流", "分享", "约", "聚会", "见面", "一起"},
		IntentEmotional: {"感觉", "觉得", "心情", "情绪", "难过", "开心", "高兴", "生气", "担心", "害怕", "紧张", "焦虑"},
		IntentRequest: {"可以", "能", "能不能", "麻烦你", "请", "求", "希望", "想要", "需要", "要求"},
		IntentOffer: {"提供", "给", "送给", "帮助", "支持", "协助", "陪同", "一起"},
		IntentComplaint: {"抱怨", "不满", "批评", "指责", "问题", "错误", "糟糕", "差劲", "不好"},
		IntentGratitude: {"谢谢", "感谢", "多谢", "感激", "辛苦", "费心了", "太好了", "棒", "赞"},
		IntentGreeting: {"你好", "大家好", "早上好", "晚上好", "嗨", "嘿", "在吗", "在忙"},
		IntentFarewell: {"再见", "拜拜", "晚安", "回见", "下次聊", "保重", "bye"},
	}
}

// initializeActionVerbs 初始化行动词
func initializeActionVerbs() []string {
	return []string{
		"做", "去", "来", "看", "听", "说", "写", "读", "买", "卖", "吃", "喝", "玩", "用", "找", "见", "约", "帮", "教", "学", "试",
		"安排", "处理", "解决", "完成", "开始", "结束", "继续", "停止", "等待", "准备", "计划", "组织", "协调", "沟通",
		"联系", "见面", "见面", "打电话", "发消息", "回复", "确认", "同意", "拒绝", "接受", "拒绝", "批准",
		"创建", "建立", "开发", "设计", "制作", "生产", "购买", "销售", "推广", "宣传", "发布",
	}
}