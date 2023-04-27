package constants

const (
	RegexSkillshareClassUrl = `https?://(?:www\.)?skillshare\.com/(?:(\w+)/)?classes/([\w-]+)/(\d{9,10})/`
	RegexSkillshareClassId  = `^\d{9,10}$`
)
