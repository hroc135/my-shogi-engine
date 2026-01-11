package model

type USICommand string

const (
	USICommandUnknown      USICommand = USICommand("")
	USICommandUSI          USICommand = USICommand("usi")
	USICommandIsReady      USICommand = USICommand("isready")
	USICommandSetOption    USICommand = USICommand("setoption")
	USICommandUSINewGame   USICommand = USICommand("usinewgame")
	USICommandPosition     USICommand = USICommand("position")
	USICommandGo           USICommand = USICommand("go")
	USICommandStop         USICommand = USICommand("stop")
	USICommandPonderhit    USICommand = USICommand("ponder")
	USICommandQuit         USICommand = USICommand("quit")
	USICommandGameover     USICommand = USICommand("gameover")
	USICommandBestMove     USICommand = USICommand("bestmove")
	USICommandGenerateMove USICommand = USICommand("generatemove")
)

func (c USICommand) String() string {
	return string(c)
}

type USIPositionSubcommand string

const (
	USIPositionSubcommandUnknown  USIPositionSubcommand = ""
	USIPositionSubcommandSFEN     USIPositionSubcommand = "sfen"
	USIPositionSubcommandStartpos USIPositionSubcommand = "startpos"
)

func (c USIPositionSubcommand) ToString() string {
	return string(c)
}

type USIGameoverSubcommand string

const (
	USIGameoverSubcommandUnknown USIGameoverSubcommand = ""
	USIGameoverSubcommandWin     USIGameoverSubcommand = "win"
	USIGameoverSubcommandLose    USIGameoverSubcommand = "lose"
	USIGameoverSubcommandDraw    USIGameoverSubcommand = "draw"
)

func (c USIGameoverSubcommand) ToString() string {
	return string(c)
}

type USIBestmoveSubcommand string

const (
	USIBestmoveSubcommandUnknown USIBestmoveSubcommand = ""
	USIBestmoveSubcommandResign  USIBestmoveSubcommand = "resign"
	USIBestmoveSubcommandWin     USIBestmoveSubcommand = "win"
	USIBestmoveSubcommandNone    USIBestmoveSubcommand = "none"
)

func (c USIBestmoveSubcommand) ToString() string {
	return string(c)
}
