package constants

import "time"

const VideoDownloadDirPath = "./output"

const EnableForwardModeButtonCallbackQuery = "button_enable_forward_mode"
const DisableForwardModeButtonCallbackQuery = "button_disable_forward_mode"
const UsePrevForwardChatButtonCallbackQuery = "button_use_prev_forward_chat"
const ChangeForwardChatButtonCallbackQuery = "button_change_forward_chat"

const TelegramMaxCaptionLen = 1024

const NonAuthSessionTTL = 10 * time.Minute
