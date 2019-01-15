package kutil


func CheckStringAlphabetOnly( value string ) ( valid bool ) {

	valid = true

	for _, c := range value {
		if false == ( c >= 0x0041 && c <= 0x005A ) && false == ( c >= 0x0061 && c <= 0x007A ) { //English
			valid = false
			break
		}
	}

	return
}

func CheckStringNumberOnly( value string ) ( valid bool ) {

	valid = true

	for _, c := range value {
		if c < 0x0030 || c > 0x0039 {
			valid = false
			break
		}
	}

	return
}

func CheckStringKoreanOnly( value string ) ( valid bool ) {

	valid = true

	for _, c := range value {
		if (c >= 0x1100 && c <= 0x1159) || (c >= 0x1161 && c <= 0x11A2) || (c >= 0x11A8 && c <= 0x11F9) ||
			(c >= 0x3131 && c <= 0x3163) || (c >= 0x3165 && c <= 0x318E) || (c >= 0xAC00 && c <= 0xD7A3) {

		} else {
			valid = false
			break
		}
	}

	return
}