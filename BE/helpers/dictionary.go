package helpers

import (
	"fmt"
	"strings"
)

func addPlusToText(text string) string {
	wordsToModify := map[string]string{
		"adalah": "ada lah", "membahas": "me bahas", "pentingnya": "penting nya",
		"semakin": "se makin", "terhubung": "ter hubung",
		"perkenalkan": "kenal kan", "halo": "hai", "kemampuan": "ke mampu an",
		"memahami": "paham", "menggunakan": "me guna kan",
		"keterampilan": "ke terampil an", "sehingga": "se hingga",
		"mengelola": "me olah", "keamanan": "ke aman an",
		"berkomunikasi": "ber komunikasi",
		"kali":          "saat", "bermain": "ber main", "membangun": "me bangun",
		"seandainya": "se andai nya", "ketinggalan": "ke tinggal an",
		"pemrograman": "pe program an", "memiliki": "me milik -i",
		"dimiliki": "di milik -i", "dijalankan": "di jalan kan",
		"pewarisan": "pe waris an", "mewarisi": "me waris -i",
		"melindungi": "me lindung -i", "membatasi": "me batas -i",
		"kelangsungan": "ke langsung an", "memungkinkan": "me mungkin kan",
		"diolah": "di olah", "berbeda": "ber beda", "berdasarkan": "ber dasar kan",
		"kelasnya": "kelas nya", "secara": "se cara",
	}

	addPlus := func(word string) string {
		if modifiedWord, exists := wordsToModify[word]; exists {
			return modifiedWord
		}
		return word
	}

	var processedText []string
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		var processedLine []string
		words := strings.Fields(line)
		for _, word := range words {
			processedLine = append(processedLine, addPlus(word))
		}
		processedText = append(processedText, strings.Join(processedLine, " "))
	}

	return strings.Join(processedText, "\n")
}

func Dictionary(str string) []string {
	str = strings.ToLower(str)
	checkWord := addPlusToText(str)

	fmt.Println(checkWord)

	dic := map[string]struct{}{
		"-i": {}, "a": {}, "ada": {}, "akan": {}, "aku": {}, "aman": {},
		"an": {}, "anak": {}, "andai": {}, "apa": {}, "atur": {}, "b": {},
		"baca": {}, "badan": {}, "bagai": {}, "bagaimana": {}, "bagi": {},
		"bagus": {}, "bahas": {}, "bahaya": {}, "bahkan": {}, "bahwa": {},
		"bangun": {}, "batas": {}, "beda": {}, "benar": {}, "ber": {},
		"besok": {}, "bicara": {}, "bisa": {}, "buah": {}, "buat": {}, "c": {},
		"cakup": {}, "cara": {}, "cari": {}, "d": {}, "dalam": {}, "dan": {},
		"dapat": {}, "dari": {}, "dasar": {}, "data": {}, "dengan": {},
		"di": {}, "dunia": {}, "e": {}, "efektif": {}, "f": {}, "fungsi": {},
		"g": {}, "guna": {}, "h": {}, "hai": {}, "hal": {}, "hari": {},
		"harus": {}, "hingga": {}, "hubung": {}, "i": {}, "ialah": {},
		"individu": {}, "induk": {}, "informasi": {}, "ingin": {}, "ini": {},
		"instansi": {}, "itu": {}, "j": {}, "jadi": {}, "jalan": {}, "k": {},
		"kah": {}, "kalau": {}, "kali": {}, "kami": {}, "kan": {}, "ke": {},
		"kelas": {}, "kemarin": {}, "kenal": {}, "kerja": {}, "kita": {},
		"l": {}, "lah": {}, "langsung": {}, "lebih": {}, "lihat": {},
		"lindung": {}, "liput": {}, "m": {}, "main": {}, "makan": {},
		"makin": {}, "mampu": {}, "mana": {}, "mau": {}, "me": {}, "metode": {},
		"milik": {}, "mudah": {}, "mulai": {}, "mungkin": {}, "n": {},
		"nama": {}, "nya": {}, "o": {}, "objek": {}, "olah": {}, "oleh": {},
		"p": {}, "paham": {}, "pe": {}, "penting": {}, "pergi": {},
		"pesawat": {}, "pribadi": {}, "program": {}, "pulang": {}, "pun": {},
		"q": {}, "r": {}, "rumah": {}, "s": {}, "saat": {}, "sarjana": {},
		"saya": {}, "se": {}, "selamat": {}, "sementara": {}, "sifat": {},
		"sini": {}, "sosial": {}, "t": {}, "teknologi": {}, "tentang": {},
		"ter": {}, "terampil": {}, "tidak": {}, "tidur": {}, "tinggal": {},
		"tulis": {}, "u": {}, "untuk": {}, "v": {}, "video": {}, "w": {},
		"waris": {}, "x": {}, "y": {}, "yang": {}, "z": {}}
	words := strings.Fields(checkWord)
	var result []string

	for _, w := range words {
		if _, found := dic[w]; found {
			result = append(result, w)
		} else {
			for _, char := range w {
				result = append(result, string(char))
			}
		}
	}

	return result
}
