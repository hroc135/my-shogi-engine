package model

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/hroc135/my-shogi-engine/util"
)

// 手番
type Color int

const (
	Black Color = iota // 先手
	White              // 後手
	NumColors
)

func (c Color) ToString() string {
	switch c {
	case Black:
		return "先手"
	case White:
		return "後手"
	default:
		return "その他"
	}
}

type Piece int

const (
	NoPiece Piece = iota // 駒なし

	// 先手の駒
	BlackPawn
	BlackLance
	BlackKnight
	BlackSilver
	BlackGold
	BlackBishop
	BlackRook
	BlackKing

	// 先手の成り駒
	BlackPromotedPawn
	BlackPromotedLance
	BlackPromotedKnight
	BlackPromotedSilver
	BlackHorse
	BlackDragon

	// 後手の駒
	WhitePawn
	WhiteLance
	WhiteKnight
	WhiteSilver
	WhiteGold
	WhiteBishop
	WhiteRook
	WhiteKing

	// 後手の成り駒
	WhitePromotedPawn
	WhitePromotedLance
	WhitePromotedKnight
	WhitePromotedSilver
	WhiteHorse
	WhiteDragon

	NumPieces
)

// PieceToString は駒の種類に対応する文字列を保持する配列
// 半角+全角+半角
var pieceToString = [...]string{
	" 　 ", " 歩 ", " 香 ", " 桂 ", " 銀 ", " 金 ", " 角 ", " 飛 ", " 王 ",
	" と ", " 杏 ", " 圭 ", " 全 ", " 馬 ", " 龍 ",
	" 歩↓", " 香↓", " 桂↓", " 銀↓", " 金↓", " 角↓", " 飛↓", " 王↓",
	" と↓", " 杏↓", " 圭↓", " 全↓", " 馬↓", " 龍↓",
}

// charToPiece 将棋の駒の文字列表現を Piece 型にマッピングする
var CharToPiece = map[byte]Piece{
	'K': BlackKing,
	'k': WhiteKing,
	'R': BlackRook,
	'r': WhiteRook,
	'B': BlackBishop,
	'b': WhiteBishop,
	'G': BlackGold,
	'g': WhiteGold,
	'S': BlackSilver,
	's': WhiteSilver,
	'N': BlackKnight,
	'n': WhiteKnight,
	'L': BlackLance,
	'l': WhiteLance,
	'P': BlackPawn,
	'p': WhitePawn,
}

var PieceToChar = map[Piece]byte{
	BlackKing:   'K',
	WhiteKing:   'k',
	BlackRook:   'R',
	WhiteRook:   'r',
	BlackBishop: 'B',
	WhiteBishop: 'b',
	BlackGold:   'G',
	WhiteGold:   'g',
	BlackSilver: 'S',
	WhiteSilver: 's',
	BlackKnight: 'N',
	WhiteKnight: 'n',
	BlackLance:  'L',
	WhiteLance:  'l',
	BlackPawn:   'P',
	WhitePawn:   'p',
}

var pieceToOpponentHandPieces = []Piece{
	NoPiece,
	WhitePawn,
	WhiteLance,
	WhiteKnight,
	WhiteSilver,
	WhiteGold,
	WhiteBishop,
	WhiteRook,
	NoPiece,
	WhitePawn,
	WhiteLance,
	WhiteKnight,
	WhiteSilver,
	WhiteBishop,
	WhiteRook,
	BlackPawn,
	BlackLance,
	BlackKnight,
	BlackSilver,
	BlackGold,
	BlackBishop,
	BlackRook,
	NoPiece,
	BlackPawn,
	BlackLance,
	BlackKnight,
	BlackSilver,
	BlackBishop,
	BlackRook,
	NumPieces,
}

func (p Piece) AsOpponentHandPiece() Piece {
	return pieceToOpponentHandPieces[p]
}

func (p Piece) ToUSIChar() byte {
	if char, ok := PieceToChar[p]; !ok {
		return 0
	} else {
		return char
	}
}

// nonPromotedToPromoted 元の駒から成り駒へのマッピング
var nonPromotedToPromoted = [NumPieces]Piece{
	BlackPawn:   BlackPromotedPawn,
	BlackLance:  BlackPromotedLance,
	BlackKnight: BlackPromotedKnight,
	BlackSilver: BlackPromotedSilver,
	BlackBishop: BlackHorse,
	BlackRook:   BlackDragon,
	WhitePawn:   WhitePromotedPawn,
	WhiteLance:  WhitePromotedLance,
	WhiteKnight: WhitePromotedKnight,
	WhiteSilver: WhitePromotedSilver,
	WhiteBishop: WhiteHorse,
	WhiteRook:   WhiteDragon,
}

func (p Piece) CanPromote() bool {
	return nonPromotedToPromoted[p] != NoPiece
}

func (p Piece) ToHumanReadableString() string {
	return pieceToString[p]
}

func (p Piece) ToHandPieceRune() rune {
	str := p.ToHumanReadableString()
	str = strings.Trim(str, " ")
	return []rune(str)[0]
}

func (p Piece) Color() Color {
	switch {
	case BlackPawn <= p && p <= BlackDragon:
		return Black
	case WhitePawn <= p && p <= WhiteDragon:
		return White
	default:
		panic(fmt.Sprintf("Piece.Color() failed with Piece: %d", int(p)))
	}
}

const BoardSize = 9

type Position struct {
	SideToMove Color // 手番
	Board      [BoardSize][BoardSize]Piece
	HandPieces [NumPieces]int // 持ち駒
	Play       int            // 手数
}

func (p Position) BoardToString() string {
	const rowSeperator = "+----+----+----+----+----+----+----+----+----+"
	var sb strings.Builder
	sb.WriteString(rowSeperator)
	sb.WriteByte('\n')
	for rank := range BoardSize {
		sb.WriteByte('|')
		for file := BoardSize - 1; file >= 0; file-- {
			sb.WriteString(p.Board[file][rank].ToHumanReadableString())
			sb.WriteByte('|')
		}
		sb.WriteByte('\n')
		sb.WriteString(rowSeperator)
		sb.WriteByte('\n')
	}

	sb.WriteString("先手 持ち駒: ")
	for blackPiece := BlackPawn; blackPiece <= BlackDragon; blackPiece++ {
		for range p.HandPieces[blackPiece] {
			sb.WriteRune(blackPiece.ToHandPieceRune())
		}
	}
	sb.WriteByte('\n')

	sb.WriteString("後手 持ち駒: ")
	for whitePiece := WhitePawn; whitePiece <= WhiteDragon; whitePiece++ {
		for range p.HandPieces[whitePiece] {
			sb.WriteRune(whitePiece.ToHandPieceRune())
		}
	}

	sb.WriteByte('\n')
	sb.WriteString("手番: ")
	switch p.SideToMove {
	case Black:
		sb.WriteString("先手")
	case White:
		sb.WriteString("後手")
	}
	sb.WriteByte('\n')

	return sb.String()
}

func (p Position) Generate(moves chan Move) {
	nonCapturePromotionMoves := make([]Move, 0)
	nonCaptureNonPromotionMoves := make([]Move, 0)

	for fileFrom := range BoardSize {
		for rankFrom := range BoardSize {
			pieceFrom := p.Board[fileFrom][rankFrom]
			if pieceFrom == NoPiece {
				continue
			}
			if pieceFrom.Color() != p.SideToMove {
				continue
			}
			for _, md := range MoveDirections[pieceFrom] {
				var distance int
				switch md.Long {
				case true:
					distance = 8
				case false:
					distance = 1
				}
				fileTo := fileFrom
				rankTo := rankFrom
				for range distance {
					fileTo += md.DeltaFile
					rankTo += md.DeltaRank
					if OutOfBoard(fileTo, rankTo) {
						continue
					}
					pieceTo := p.Board[fileTo][rankTo]
					if pieceTo != NoPiece && pieceTo.Color() == p.SideToMove {
						// 味方の駒がいる
						break
					}

					if pieceFrom.CanPromote() && promotablePosition(p.SideToMove, rankFrom, rankTo) {
						// 成って移動する
						move := Move{
							FileFrom:   fileFrom,
							RankFrom:   rankFrom,
							PieceFrom:  pieceFrom,
							FileTo:     fileTo,
							RankTo:     rankTo,
							PieceTo:    pieceTo,
							Promotion:  true,
							SideToMove: p.SideToMove,
						}
						if pieceTo != NoPiece && pieceTo.Color() != p.SideToMove {
							// 駒を取る指し手の場合はすぐに返す
							move.Drop = true
							moves <- move
							break
						} else {
							// 駒を取らない指し手の場合は後で返す
							nonCapturePromotionMoves = append(nonCapturePromotionMoves, move)
						}
					}

					if CanPutWithoutPromotion(pieceFrom, rankTo) {
						// 成らずに移動する
						move := Move{
							FileFrom:   fileFrom,
							RankFrom:   rankFrom,
							PieceFrom:  pieceFrom,
							FileTo:     fileTo,
							RankTo:     rankTo,
							PieceTo:    pieceTo,
							Promotion:  false,
							SideToMove: p.SideToMove,
						}
						if pieceTo != NoPiece && pieceTo.Color() != p.SideToMove {
							// 駒を取る指し手の場合はすぐに返す
							move.Drop = true
							moves <- move
						} else {
							// 駒を取らない指し手の場合は後で返す
							nonCaptureNonPromotionMoves = append(nonCaptureNonPromotionMoves, move)
						}

					}

					if pieceTo != NoPiece {
						// 相手の駒なので、ここで利きが止まる
						break
					}
				}
			}
		}
	}

	// 駒を取らない、成る指し手
	for _, m := range nonCapturePromotionMoves {
		moves <- m
	}
	// 駒を取らない、成らない指し手
	for _, m := range nonCaptureNonPromotionMoves {
		moves <- m
	}

	// 駒を打つ指し手
	var minPiece, maxPiece Piece
	switch p.SideToMove {
	case White:
		minPiece = WhitePawn
		maxPiece = WhiteRook
	case Black:
		minPiece = BlackPawn
		maxPiece = BlackRook
	}
	for pieceFrom := minPiece; pieceFrom < maxPiece; pieceFrom++ {
		if p.HandPieces[pieceFrom] == 0 {
			continue
		}
		for fileTo := range BoardSize {
			for rankTo := range BoardSize {
				if p.Board[fileTo][rankTo] != NoPiece {
					continue
				}
				if !CanPutWithoutPromotion(pieceFrom, rankTo) {
					continue
				}
				if (pieceFrom == WhitePawn || pieceFrom == BlackPawn) && p.IsPawnExist(fileTo) {
					// 二歩
					continue
				}
				move := Move{
					FileFrom:   -1,
					RankFrom:   -1,
					PieceFrom:  pieceFrom,
					FileTo:     fileTo,
					RankTo:     rankTo,
					PieceTo:    pieceFrom,
					Drop:       false,
					Promotion:  false,
					SideToMove: p.SideToMove,
				}
				moves <- move
			}
		}
	}
}

func (p Position) IsPawnExist(file int) bool {
	switch p.SideToMove {
	case White:
		for _, piece := range p.Board[file] {
			if piece == WhitePawn {
				return true
			}
		}
	default:
		for _, piece := range p.Board[file] {
			if piece == BlackPawn {
				return true
			}
		}
	}
	return false
}

func promotablePosition(side Color, rankFrom, rankTo int) bool {
	return (side == Black && rankTo <= 2) ||
		(side == White && rankTo >= 6) ||
		(side == Black && rankFrom <= 2) ||
		(side == White && rankFrom >= 6)
}

func CanPutWithoutPromotion(pieceFrom Piece, rankTo int) bool {
	return (pieceFrom == BlackPawn && rankTo >= 1) ||
		(pieceFrom == BlackLance && rankTo >= 1) ||
		(pieceFrom == BlackKnight && rankTo >= 2) ||
		(pieceFrom == WhitePawn && rankTo <= 7) ||
		(pieceFrom == WhiteLance && rankTo <= 7) ||
		(pieceFrom == WhiteKnight && rankTo <= 6) ||
		(pieceFrom != BlackPawn && pieceFrom != BlackLance && pieceFrom != BlackKnight && pieceFrom != WhitePawn && pieceFrom != WhiteLance && pieceFrom != WhiteKnight)
}

const StartPosSFEN string = "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"

func SetPosition(sfen string) Position {
	p := Position{}

	file := BoardSize - 1
	rank := 0
	promotion := false

	i := 0
	for ; ; i++ {
		char := sfen[i]
		if char == ' ' {
			break
		}
		// 次の段へ
		if char == '/' {
			rank++
			file = BoardSize - 1
			continue
		}
		// 成り駒
		if char == '+' {
			promotion = true
			continue
		}
		// 数字の分だけ空白マス
		if unicode.IsDigit(rune(char)) {
			empty := char - '0'
			for range empty {
				p.Board[file][rank] = NoPiece
				file--
			}
			continue
		}
		// 駒を盤面に配置
		piece := CharToPiece[char]
		if promotion {
			p.Board[file][rank] = nonPromotedToPromoted[piece]
		} else {
			p.Board[file][rank] = piece
		}
		file--
		promotion = false
	}
	i++

	// 手番
	if sfen[i] == 'b' {
		p.SideToMove = Black
	} else {
		p.SideToMove = White
	}
	i += 2

	// 持ち駒
	cnt := 0
	for ; ; i++ {
		char := sfen[i]
		if char == ' ' {
			break
		}
		if char == '-' {
			continue
		}
		if unicode.IsDigit(rune(char)) {
			cnt = cnt*10 + int(char-'0')
			continue
		}
		p.HandPieces[CharToPiece[char]] = max(1, cnt)
		cnt = 0
	}
	i++

	// 手数
	play, err := strconv.Atoi(sfen[i:])
	if err != nil {
		fmt.Printf("invalid play number: %s\n", sfen[i:])
	}
	p.Play = play

	return p
}

func OutOfBoard(f, r int) bool {
	return f < 0 || BoardSize <= f || r < 0 || BoardSize <= r
}

type Move struct {
	FileFrom   int
	RankFrom   int
	PieceFrom  Piece
	FileTo     int
	RankTo     int
	PieceTo    Piece
	Drop       bool // 駒を取った
	Promotion  bool
	SideToMove Color
}

var (
	// Resign 投了
	Resign = Move{
		FileFrom: 2,
		FileTo:   2,
	}

	// Win 入玉勝ち
	Win = Move{
		FileFrom: 3,
		FileTo:   3,
	}

	// None 無効な指し手
	None = Move{
		FileFrom: 4,
		FileTo:   4,
	}

	// FileToZenkakuLetters = [...]string{"１", "２", "３", "４", "５", "６", "７", "８", "９"}
	RankToKanjiLetters = [...]string{"一", "二", "三", "四", "五", "六", "七", "八", "九"}
)

func (m Move) ToString() string {
	var side string
	switch m.SideToMove {
	case Black:
		side = "△"
	case White:
		side = "▲"
	}
	var promotionPrefix string
	if m.Promotion {
		promotionPrefix = "成"
	}
	return fmt.Sprintf("%s%d%s%s%s",
		side,
		m.FileTo+1,
		RankToKanjiLetters[m.RankTo],
		strings.Trim(m.PieceFrom.ToHumanReadableString(), " "),
		promotionPrefix,
	)
}

func (m Move) Equals(move Move) bool {
	return reflect.DeepEqual(m, move)
}

func (m Move) GetHashCode() int {
	seed := 0
	seed ^= m.FileFrom + (seed << 6) + (seed >> 2)
	seed ^= m.RankFrom + (seed << 6) + (seed >> 2)
	seed ^= int(m.PieceFrom) + (seed << 6) + (seed >> 2)
	seed ^= m.FileTo + (seed << 6) + (seed >> 2)
	seed ^= m.RankTo + (seed << 6) + (seed >> 2)
	seed ^= int(m.PieceTo) + (seed << 6) + (seed >> 2)
	seed ^= util.BoolToInt(m.Drop) + (seed << 6) + (seed >> 2)
	seed ^= util.BoolToInt(m.Promotion) + (seed << 6) + (seed >> 2)
	seed ^= int(m.SideToMove) + (seed << 6) + (seed >> 2)
	return seed
}

func (m Move) ToUSIString() string {
	if m == Resign {
		return USIBestmoveSubcommandResign.ToString()
	}
	if m == Win {
		return USIBestmoveSubcommandWin.ToString()
	}

	var usiStr string
	if m.Drop {
		usiStr += strings.ToUpper(string(m.PieceFrom.ToUSIChar()))
		usiStr += "*"
	} else {
		usiStr += fmt.Sprint(m.FileFrom + '1')
		usiStr += fmt.Sprint(m.RankFrom + 'a')
	}
	usiStr += fmt.Sprint(m.FileTo + '1')
	usiStr += fmt.Sprint(m.RankTo + 'a')

	if m.Promotion {
		usiStr += "+"
	}

	return usiStr
}

func (m Move) FromUSIString(pos Position, moveStr string) Move {
	switch USIBestmoveSubcommand(moveStr) {
	case USIBestmoveSubcommandResign:
		return Resign
	case USIBestmoveSubcommandWin:
		return Win
	case USIBestmoveSubcommandNone:
		return None
	}

	move := Move{}
	if moveStr[1] == '*' { // 駒打ちの指し手
		move.FileFrom = -1
		move.RankFrom = -1
		move.PieceFrom = CharToPiece[moveStr[0]]
		if pos.SideToMove == White {
			move.PieceFrom = move.PieceFrom.AsOpponentHandPiece()
		}
		move.Drop = true
	} else {
		move.FileFrom = int(moveStr[0] - '1')
		move.RankFrom = int(moveStr[1] - 'a')
		move.PieceFrom = pos.Board[move.FileFrom][move.RankFrom]
	}
	move.FileTo = int(moveStr[2] - '1')
	move.RankTo = int(moveStr[3] - 'a')
	move.PieceTo = pos.Board[move.FileTo][move.RankTo]
	if len(moveStr) == 5 {
		move.Promotion = true
	}
	move.SideToMove = pos.SideToMove
	return move
}

type Direction struct {
	DeltaFile int
	DeltaRank int
}

var (
	upLeft    = Direction{DeltaFile: +1, DeltaRank: -1}
	up        = Direction{DeltaFile: 0, DeltaRank: -1}
	upRight   = Direction{DeltaFile: -1, DeltaRank: -1}
	left      = Direction{DeltaFile: +1, DeltaRank: 0}
	right     = Direction{DeltaFile: -1, DeltaRank: 0}
	downLeft  = Direction{DeltaFile: +1, DeltaRank: +1}
	down      = Direction{DeltaFile: 0, DeltaRank: +1}
	downRight = Direction{DeltaFile: -1, DeltaRank: +1}
)

type MoveDirection struct {
	Direction
	Long bool
}

var MoveDirections = [][]MoveDirection{
	// NoPiece
	{},
	// BlackPawn
	{
		{Direction: up},
	},
	// BlackLance
	{
		{Direction: up, Long: true},
	},
	// BlackKnight
	{
		{Direction: Direction{DeltaFile: +1, DeltaRank: -2}},
		{Direction: Direction{DeltaFile: -1, DeltaRank: -2}},
	},
	// BlackSilver
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: downLeft},
		{Direction: downRight},
	},
	// BlackGold
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: left},
		{Direction: right},
		{Direction: down},
	},
	// BlackBishop
	{
		{Direction: upLeft, Long: true},
		{Direction: upRight, Long: true},
		{Direction: downLeft, Long: true},
		{Direction: downRight, Long: true},
	},
	// BlackRook
	{
		{Direction: up, Long: true},
		{Direction: left, Long: true},
		{Direction: right, Long: true},
		{Direction: down, Long: true},
	},
	// BlackKing
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: left},
		{Direction: right},
		{Direction: downLeft},
		{Direction: down},
		{Direction: downRight},
	},
	// BlackPromotedPawn
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: left},
		{Direction: right},
		{Direction: down},
	},
	// BlackPromotedLance
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: left},
		{Direction: right},
		{Direction: down},
	},
	// BlackPromotedKnight
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: left},
		{Direction: right},
		{Direction: down},
	},
	// BlackPromotedSilver
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: left},
		{Direction: right},
		{Direction: down},
	},
	// BlackHorse
	{
		{Direction: upLeft, Long: true},
		{Direction: up},
		{Direction: upRight, Long: true},
		{Direction: left},
		{Direction: right},
		{Direction: downLeft, Long: true},
		{Direction: down},
		{Direction: downRight, Long: true},
	},
	// BlackDragon
	{
		{Direction: upLeft},
		{Direction: up, Long: true},
		{Direction: upRight},
		{Direction: left, Long: true},
		{Direction: right, Long: true},
		{Direction: downLeft},
		{Direction: down, Long: true},
		{Direction: downRight},
	},
	// WhitePawn
	{
		{Direction: down},
	},
	// WhiteLance
	{
		{Direction: down, Long: true},
	},
	// WhiteKnight
	{
		{Direction: Direction{DeltaFile: +1, DeltaRank: +2}},
		{Direction: Direction{DeltaFile: -1, DeltaRank: +2}},
	},
	// WhiteSilver
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: downLeft},
		{Direction: downRight},
	},
	// WhiteGold
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: left},
		{Direction: right},
		{Direction: down},
	},
	// WhiteBishop
	{
		{Direction: upLeft, Long: true},
		{Direction: upRight, Long: true},
		{Direction: downLeft, Long: true},
		{Direction: downRight, Long: true},
	},
	// WhiteRook
	{
		{Direction: up, Long: true},
		{Direction: left, Long: true},
		{Direction: right, Long: true},
		{Direction: down, Long: true},
	},
	// WhiteKing
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: left},
		{Direction: right},
		{Direction: downLeft},
		{Direction: down},
		{Direction: downRight},
	},
	// WhitePromotedPawn
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: left},
		{Direction: right},
		{Direction: down},
	},
	// WhitePromotedLance
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: left},
		{Direction: right},
		{Direction: down},
	},
	// WhitePromotedKnight
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: left},
		{Direction: right},
		{Direction: down},
	},
	// WhitePromotedSilver
	{
		{Direction: upLeft},
		{Direction: up},
		{Direction: upRight},
		{Direction: left},
		{Direction: right},
		{Direction: down},
	},
	// WhiteHorse
	{
		{Direction: upLeft, Long: true},
		{Direction: up},
		{Direction: upRight, Long: true},
		{Direction: left},
		{Direction: right},
		{Direction: downLeft, Long: true},
		{Direction: down},
		{Direction: downRight, Long: true},
	},
	// WhiteDragon
	{
		{Direction: upLeft},
		{Direction: up, Long: true},
		{Direction: upRight},
		{Direction: left, Long: true},
		{Direction: right, Long: true},
		{Direction: downLeft},
		{Direction: down, Long: true},
		{Direction: downRight},
	},
	{},
}
