package msgs

const (
	IconThumbsUp      = "👍"
	IconX             = "❌"
	IconDevil         = "😈"
	IconPointingRight = "👉"
	IconPointingDown  = "👇"
	IconSkull         = "💀"
	IconWarning       = "⚠"
	IconAlarmClock    = "⏰"
	IconPrevious      = "⇦"
	IconNext          = "⇨"

	MsgThumbsUp       = IconThumbsUp
	MsgCantUnderstand = IconX + " -- Desculpe, não entendi"
	MsgNotAuthorized  = IconDevil + " -- Desculpe, você não tem permissão para isso"
	MsgServerError    = IconSkull + " -- Desculpe, tem algo de errado comigo..."
	MsgNotFound       = IconWarning + " -- Desculpe, não consegui encontrar isso"
	MsgHelp           = IconThumbsUp + " -- Os comandos são:\n/admiral\n/anime\n/manga\n/money\n/movie"
	MsgAdmiral        = IconWarning + " -- The Admiral command is /admiral <admiral name> "
	MsgAnime          = IconWarning + " -- O comando é /anime <nome do anime>\nO resultado é baseado em uma pesquisa no MyanimeList"
	MsgManga          = IconWarning + " -- O comando é /manga <nome do mangá>\nO resultado é baseado em uma pesquisa no MyanimeList"
	MsgMoney          = IconWarning + "-- O comando é /money <quantidade> <moeda principal> <moeda a ser convertida>"
	MsgMovie          = IconWarning + "-- O comando é /movie <nome do filme> O resultado é baseado em uma pesquisa do MovieDB"
)
