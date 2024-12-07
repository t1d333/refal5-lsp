package symbols

import "unicode/utf8"

func ByteOffsetToRunePosition(text string, byteOffset int) int {
	runes := []rune(text)
	runePosition := 0
	currentByteOffset := 0

	for _, r := range runes {
		if currentByteOffset >= byteOffset {
			break
		}
		currentByteOffset += utf8.RuneLen(r)
		runePosition++
	}

	return runePosition
}

func RunePositionToByteOffset(text string, runePosition int) int {
	runes := []rune(text)
	if runePosition > len(runes) {
		return -1 // Позиция за пределами строки
	}

	return len([]byte(string(runes[:runePosition])))
}

func RuneEndPosition(text string, runePosition int) int {
	runes := []rune(text)
	if runePosition > len(runes) {
		return -1 // Позиция за пределами строки
	}

	return len([]byte(string(runes[:runePosition]))) + utf8.RuneLen(runes[runePosition]) - 1
}
