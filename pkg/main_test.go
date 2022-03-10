package pkg

var (
	ExportedEnvFilePath            = envFilePath
	ExportedReadEnv                = readEnv
	ExportedNewConfig              = newConfig
	ExportedNewIFTTT               = newIFTTT
	ExportedNewSlackClient         = newSlackClient
	ExportedGetConversationHistory = (*slackClient).getConversationHistory
	ExportedFetchMessages          = (*slackClient).fetchMessages
	ExportedFilterMessages         = (*slackClient).filterMessages
	ExportedSortMessages           = (*slackClient).sortMessages
)

type (
	ExportedIFTTT                   = ifttt
	ExportedSlackMessage            = slackMessage
	ExportedFilterSlackMessagesArgs = filterSlackMessagesArgs
	ExportedSlackClient             = slackClient
)
