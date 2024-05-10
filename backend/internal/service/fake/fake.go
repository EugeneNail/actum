package fake

import (
	"math/rand/v2"
	"strings"
)

var lorem []string
var names []string
var addresses []string
var servers []string
var domains []string

const COMMA_CHANCE = 0.1

func init() {
	lorem = []string{
		"sed", "in", "et", "ut", "ac", "sit", "amet", "quis", "id", "nec", "eget", "vitae", "eu", "at", "a", "vel", "nulla", "non", "nunc", "tincidunt", "aliquam", "ipsum", "pellentesque", "mauris", "orci", "turpis", "erat",
		"ante", "donec", "vestibulum", "nisi", "tellus", "purus", "elit", "diam", "sapien", "odio", "lectus", "est", "ligula", "dui", "neque", "arcu", "malesuada", "libero", "quam", "lorem", "risus", "porta", "justo", "felis",
		"leo", "egestas", "augue", "ultrices", "massa", "fringilla", "enim", "auctor", "venenatis", "velit", "tristique", "viverra", "sem", "faucibus", "nibh", "eros", "metus", "fermentum", "ex", "magna", "luctus", "sollicitudin",
		"molestie", "tortor", "posuere", "nisl", "iaculis", "urna", "etiam", "vivamus", "vehicula", "feugiat", "dolor", "rutrum", "rhoncus", "condimentum", "varius", "sodales", "consectetur", "volutpat", "consequat", "blandit",
		"aliquet", "lacus", "congue", "maximus", "gravida", "eleifend", "accumsan", "ullamcorper", "semper", "pharetra", "nullam", "morbi", "maecenas", "interdum", "dignissim", "dapibus", "phasellus", "mi", "curabitur", "commodo",
		"aenean", "tempus", "lacinia", "hendrerit", "euismod", "suspendisse", "scelerisque", "proin", "nam", "elementum", "dictum", "ultricies", "mollis", "duis", "cras", "porttitor", "placerat", "fusce", "cursus", "lobortis",
		"convallis", "tempor", "pretium", "suscipit", "finibus", "facilisis", "bibendum", "ornare", "laoreet", "integer", "imperdiet", "mattis", "pulvinar", "sagittis", "vulputate", "quisque", "efficitur", "praesent", "fames",
		"primis", "senectus", "netus", "habitant", "curae", "cubilia", "adipiscing", "facilisi", "ridiculus", "platea", "penatibus", "parturient", "natoque", "nascetur", "mus", "montes", "magnis", "hac", "habitasse", "dis",
		"dictumst", "per", "potenti", "torquent", "taciti", "sociosqu", "nostra", "litora", "inceptos", "himenaeos", "conubia", "class", "aptent", "ad",
	}

	names = []string{
		"Alma", "Cecilia", "Lyle", "Vanessa", "Chelsea", "Jessica", "Grant", "Toby", "Regina", "Seth", "Theodore", "Bonita", "Ronald", "Maggie", "Victoria", "Nicole",
		"Nathaniel", "Annie", "Lester", "Teri", "Alejandro", "Monique", "Marion", "Laverne", "Tiffany", "Jana", "Lucy", "Oliver", "Leona", "Carol", "Herbert", "Lillie",
		"Sammy", "Rafael", "Claire", "Ben", "Mabel", "Priscilla", "Freddie", "Judy", "Leah", "Gerard", "Henry", "Margaret", "Jacquelyn", "Jeanette", "Tara", "Doug", "Jan", "Cecelia",
	}

	addresses = []string{
		"sethu.rich", "f5il5z5vuibbzo2", "finn.fleming", "uy0o6m9gjlpa19cx", "kylan.anderson", "ny6155ihvalevit", "eng.delacruz", "xh61hotavhajtp", "zayne.harvey", "wb8g8wozfz8kiuu0uj5p", "ahoua.soto",
		"n8fwcbzzmpdwwpcc76", "james-paul.baldwin", "hkn4occk1lr", "klein.patterson", "r6enwccag4zxx27q", "geordan.wyatt", "gwkmwxg5umu031h0", "ciann.robertson", "wqmlduuodqhnyhix2nbj", "chiron.stokes",
		"j4jhjmrh9", "gianluca.perez", "hti18hnb7rajhbomv", "aslam.bryan", "pcf6us8o1l7rb0", "jerrick.cannon", "ovrc9rcix0bb", "conghaile.hopper", "kqcz18xblb", "humza.stevenson", "xucm82jf0uw",
	}

	servers = []string{"@106-list", "@rediffmail", "@aol", "@outlook", "@ymail", "@comcast", "@yahoo", "@hotmail", "@googlemail", "@msn", "@gmail"}

	domains = []string{".com", ".eu", ".net", ".ru", ".org", ".cc", ".co", ".us"}

}

func getRandomItem(array []string) string {
	return array[rand.IntN(len(array))]
}

func Name() string {
	return getRandomItem(names)
}

func Email() string {
	return getRandomItem(addresses) + getRandomItem(servers) + getRandomItem(domains)
}

func Word() string {
	return getRandomItem(lorem)
}

func WordLength(min int, max int) (word string) {
	length := len(word)
	for length < min || length > max {
		word = Word()
	}
	return
}

func Sentence() string {
	return SentenceLength(5, 20)
}

func SentenceLength(min int, max int) string {
	var words []string
	length := rand.IntN(max-min) + min

	for i := 0; i < length; i++ {
		word := Word()
		if rand.Float32() < COMMA_CHANCE {
			word += ","
		}
		words = append(words, word)
	}

	words[0] = strings.ToUpper(words[0][:1]) + words[0][1:]

	return strings.Join(words, " ") + "."
}

func Paragraph() string {
	return ParagraphLength(3, 5)
}

func ParagraphLength(min int, max int) string {
	var sentences []string

	length := rand.IntN(max-min) + min

	for i := 0; i < length; i++ {
		sentences = append(sentences, Sentence())
	}

	return strings.Join(sentences, " ")
}

func Text() string {
	return TextLength(1, 5)
}

func TextLength(min int, max int) string {
	var paragraphs []string

	length := rand.IntN(max-min) + min

	for i := 0; i < length; i++ {
		paragraphs = append(paragraphs, Paragraph())
	}

	return strings.Join(paragraphs, "\n")

	return ""
}
