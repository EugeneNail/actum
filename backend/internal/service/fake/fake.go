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
var icons []string

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

	icons = []string{
		"Man", "Woman", "Wc", "Family_Restroom", "Escalator_Warning", "Hail", "Accessibility", "Accessible", "Directions_Run", "Airline_Seat_Recline_Normal", "Nature_People", "Bathtub", "Hot_Tub", "Baby_Changing_Station",
		"Supervisor_Account", "Record_Voice_Over", "Pregnant_Woman", "Groups", "People", "Engineering", "Emoji_People", "Self_Improvement", "Connect_Without_Contact", "Hiking", "Reduce_Capacity", "Diversity_3", "Elderly",
		"Personal_Injury", "Diversity_1", "Hotel", "Sports_Gymnastics", "Sports_Kabaddi", "Kayaking", "Skateboarding", "Snowshoeing", "Pool", "Rowing", "Surfing", "Sports_Handball", "Paragliding", "Downhill_Skiing",
		"Sports_Martial_Arts", "Wb_Sunny", "Nightlight_Round", "Nights_Stay", "Filter_Drama", "Thunderstorm", "Snowing", "Ac_Unit", "Wb_Twilight", "Spa", "Grass", "Park", "Local_Florist", "Forest", "Terrain", "Home",
		"Chair_Alt", "Door_Sliding", "Table_Bar", "Shelves", "Wind_Power", "Roller_Shades", "Countertops", "Bed", "Shower", "Light", "Meeting_Room", "Checkroom", "Kitchen", "Local_Laundry_Service", "Headphones", "Keyboard",
		"Tv", "Mouse", "Router", "King_Bed", "Tablet_Mac", "Pedal_Bike", "Agriculture", "Two_Wheeler", "Directions_Car", "Directions_Bus_Filled", "Train", "Tram", "Sailing", "Directions_Boat", "Flight_Takeoff", "Fire_Truck",
		"Snowmobile", "Rocket_Launch", "House", "Maps_Home_Work", "Apartment", "Location_City", "Domain", "Emoji_Transportation", "Factory", "Store_Mall_Directory", "Festival", "Church", "Castle", "Storefront", "Account_Balance",
		"Store", "Delete", "Visibility", "Favorite", "Description", "Lock", "Schedule", "Language", "Thumb_Up", "Filter_Alt", "Event", "Dashboard", "Paid", "Question_Answer", "Article", "Lightbulb", "Credit_Card", "History",
		"Trending_Up", "Fact_Check", "Account_Balance_Wallet", "Build", "Analytics", "Receipt", "Explore", "Pending_Actions", "Leaderboard", "Thumb_Up_Off_Alt", "Card_Giftcard", "View_In_Ar", "Timeline", "Stars", "Dns",
		"Space_Dashboard", "Alarm", "Bug_Report", "Gavel", "Pan_Tool", "Extension", "Hourglass_Empty", "Thumb_Down", "Support", "Loyalty", "Euro_Symbol", "Table_View", "Track_Changes", "Perm_Media", "Backup", "File_Present",
		"Trending_Down", "Percent", "Shopping_Cart", "Shopping_Bag", "Swipe", "Work", "Print", "Room", "Translate", "Book_Online", "Perm_Phone_Msg", "G_Translate", "Aspect_Ratio", "Thumbs_Up_Down", "Theaters", "Tour", "Mark_As_Unread",
		"Settings_Input_Antenna", "Balance", "View_Carousel", "All_Inbox", "Settings_Remote", "Settings_Voice", "Online_Prediction", "Camera_Enhance", "Fax", "Satellite_Alt", "Settings_Cell", "App_Blocking", "Barcode_Reader",
		"Payments", "Share", "School", "Public", "Emoji_Events", "Notifications_Active", "Construction", "Psychology", "Health_And_Safety", "Water_Drop", "Notifications_None", "Sports_Esports", "Workspace_Premium",
		"Precision_Manufacturing", "Military_Tech", "Science", "History_Edu", "Handshake", "Coronavirus", "Sports_Soccer", "Recycling", "Waving_Hand", "Luggage", "Vaccines", "Interests", "Sports_Basketball", "Sports",
		"Heart_Broken", "Sports_Tennis", "Deck", "Scale", "Sports_Motorsports", "Sanitizer", "Sports_Baseball", "Party_Mode", "Mail", "Flag", "Push_Pin", "Create", "Photo_Camera", "Image", "Tune", "Auto_Stories", "Palette",
		"Music_Note", "Healing", "Vpn_Key", "Stay_Current_Portrait", "Map", "Restaurant", "Local_Fire_Department", "Volunteer_Activism", "Celebration", "Local_Police", "Local_Gas_Station", "Electrical_Services", "Traffic",
		"Theater_Comedy", "Mic", "Volume_Up", "Imagesearch_Roller", "Memory", "Fitness_Center", "Business_Center", "Beach_Access", "Casino", "Vaping_Rooms",
	}
}

func getRandomItem(array []string) string {
	return array[rand.IntN(len(array))]
}

func Icon() string {
	return getRandomItem(icons)
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

	return strings.Join(paragraphs, "\\n")
}
