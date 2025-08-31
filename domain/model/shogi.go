package model

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

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

func AsPromoted(p Piece) (Piece, error) {
	if nonPromotedToPromoted[p] == NoPiece {
		return NoPiece, fmt.Errorf("AsPromoted: piece %s cannot promote", pieceToString[p])
	}
	return nonPromotedToPromoted[p], nil
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
			sb.WriteString(p.Board[file][rank].ToString())
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
