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
	IconPrevious      = "❰"
	IconNext          = "❯"

	MsgThumbsUp       = IconThumbsUp
	MsgCantUnderstand = IconX + " -- Desculpe, não entendi"
	MsgNotAuthorized  = IconDevil + " -- Desculpe, você não tem permissão para isso"
	MsgServerError    = IconSkull + " -- Desculpe, tem algo de errado comigo..."
	MsgNotFound       = IconWarning + " -- Desculpe, não consegui encontrar isso"
	MsgHelp           = IconThumbsUp + " -- Os comandos são:\n/admirals\n/animes\n/mangas\n/money\n/movies"
	MsgAdmirals       = IconWarning + " -- The Admiral command is /admirals <admiral name> "
	MsgAnimes         = IconWarning + " -- O comando é /animes <nome do anime>\nO resultado é baseado em uma pesquisa no MyanimeList"
	MsgMangas         = IconWarning + " -- O comando é /mangas <nome do mangá>\nO resultado é baseado em uma pesquisa no MyanimeList"
	MsgMoney          = IconWarning + "-- O comando é /money <quantidade> <moeda principal> <moeda a ser convertida>"
	MsgMovies         = IconWarning + "-- O comando é /movies <nome do filme> O resultado é baseado em uma pesquisa do MovieDB"
	MsgTVShow         = IconWarning + "-- O comando é /series <nome da serie> O resultado é baseado em uma pesquisa do MovieDB"
)
