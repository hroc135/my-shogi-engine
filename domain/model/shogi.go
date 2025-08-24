package model

import "strings"

// 手番
type Color int

const (
	Black Color = iota // 先手
	White              // 後手
	NumColors
)

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
var pieceToString = [...]string{
	"　　", " 歩 ", " 香 ", " 桂 ", " 銀 ", " 金 ", " 角 ", " 飛 ", " 王 ",
	" と ", " 杏 ", " 圭 ", " 全 ", " 馬 ", " 龍 ",
	"歩↓", "香↓", "桂↓", "銀↓", "金↓", "角↓", "飛↓", "王↓",
	"と↓", "杏↓", "圭↓", "全↓", "馬↓", "龍↓",
}

func (p Piece) ToString() string {
	return pieceToString[int(p)]
}

func (p Piece) ToHandPieceRune() rune {
	str := p.ToString()
	str = strings.Trim(str, " ")
	return []rune(str)[0]
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
			sb.WriteString(p.Board[rank][file].ToString())
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
