package contentgenerator

import (
	"bytes"
	"encoding/json"
	"jlpt-practice-helper/jsonfileflow/model"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var artists = map[string]string{
	"Drake":                    "HipHop, Trap, male vocals",
	"Bruno Mars":               "Funk, Dance Pop, Groovy, male vocals",
	"Fleetwood Mac":            "Classic Rock, Mellifluous",
	"Ed Sheeran":               "Folk, Acoustic Guitar, male vocals",
	"Tim McGraw":               "Country, Americana, male vocals",
	"Elton John":               "Piano Pop Rock, Theatrical, male vocals",
	"Dolly Parton":             "Country, Storytelling, female vocals",
	"Red Hot Chili Peppers":    "Funk Rock, Stadium, heavy drums",
	"Coldplay":                 "Alternative Rock, Atmospheric",
	"Taylor Swift":             "Pop, Alternative Folk, Emotional, female vocals",
	"Elvis Presley":            "50s Rock, Hero Theme, male vocals",
	"Adele":                    "Soul, Emotional, Torch-Lounge, female vocals",
	"Ariana Grande":            "Pop, Dance Pop, Ethereal, female vocals",
	"Billie Eilish":            "Pop, Dark, Minimal, female vocals",
	"The Weeknd":               "RnB, Dark, Cinematic, male vocals",
	"Beyoncé":                  "RnB, Anthemic, Danceable, female vocals",
	"Kendrick Lamar":           "HipHop, Lyrical, Storytelling, male vocals",
	"Lady Gaga":                "Pop, Theatrical, Dance, female vocals",
	"Jay-Z":                    "HipHop, Aggressive, Storytelling, male vocals",
	"Rihanna":                  "RnB, Dance Pop, Festive, female vocals",
	"Kanye West":               "HipHop, Progressive, Eclectic, male vocals",
	"Justin Bieber":            "Pop, Danceable, Chillwave, male vocals",
	"Katy Perry":               "Pop, Glitter, Festive, female vocals",
	"Snoop Dogg":               "Rap, Funk, Chill, male vocals",
	"Metallica":                "Heavy Metal, Power",
	"AC/DC":                    "Hard Rock, Stomp",
	"Madonna":                  "Dance Pop, High-NRG, female vocals",
	"David Bowie":              "70s British Rock, Art, Eclectic, male vocals",
	"Bob Dylan":                "Folk, Storytelling, Acoustic Guitar, male vocals",
	"Post Malone":              "Rap, Ethereal, Ambient, male vocals",
	"Maroon 5":                 "Pop Rock, Danceable, male vocals",
	"Shakira":                  "Latin, Dance Pop, Festive, female vocals",
	"Dua Lipa":                 "Disco, Dance Pop, Groovy, female vocals",
	"Michael Jackson":          "80s Pop, Dance, Iconic, male vocals",
	"Prince":                   "Funk, Eclectic, Glam, male vocals",
	"Miley Cyrus":              "Pop, Rock, Party, female vocals",
	"Cardi B":                  "Rap, Aggressive, Party, female vocals",
	"Imagine Dragons":          "Rock, Anthemic, Emotion",
	"Camila Cabello":           "Pop, Latin Jazz, Festive, female vocals",
	"Harry Styles":             "Pop, Rock, Groovy, male vocals",
	"Sam Smith":                "Soul, Emotional, Lounge, male vocals",
	"Lizzo":                    "Pop, Funk, Empowering, female vocals",
	"Daft Punk":                "Electronic, Dance, Futuristic",
	"Gorillaz":                 "Alternative Rock, Electronic, Unusual",
	"The Beatles":              "60s British Pop, Classic, Rock",
	"Queen":                    "Rock, Operatic, Theatrical, Male Vocals",
	"Led Zeppelin":             "Hard Rock, Blues Rock, Epic",
	"Pink Floyd":               "Rock, Progressive, Atmospheric",
	"The Rolling Stones":       "Rock, Blues Rock, Classic",
	"Bob Marley":               "Reggae, Peaceful, Soulful, male vocals",
	"Frank Sinatra":            "1940s big band, Lounge Singer, male vocals",
	"Aretha Franklin":          "Soul, Gospel, Powerful, female vocals",
	"Whitney Houston":          "Pop, RnB, Emotional, female vocals",
	"Stevie Wonder":            "Soul, Funk, Joyful, male vocals",
	"The Chainsmokers":         "EDM, Dance, Party",
	"Nicki Minaj":              "Rap, Danceable, Bold, female vocals",
	"Green Day":                "Punk Rock, Aggressive, Youthful",
	"Nirvana":                  "Grunge, Dark, Raw, Male Vocals",
	"Amy Winehouse":            "Soul, Jazz, Torch-Lounge, female vocals",
	"Linkin Park":              "Rock, Nu Metal, Emotional",
	"Aerosmith":                "Rock, Hard Rock, Classic",
	"Bon Jovi":                 "Rock, Anthem, Stadium",
	"Billy Joel":               "Pop, Rock, Storytelling, male vocals",
	"Phil Collins":             "Pop, Rock, Emotional, soundtrack, male vocals",
	"Genesis":                  "Rock, Progressive, Art",
	"The Eagles":               "Rock, Country Rock, Harmonious",
	"The Doors":                "Rock, Psychedelic, Mysterious",
	"Janis Joplin":             "Rock, Blues Rock, Raw Emotion, female vocals",
	"Jimi Hendrix":             "Rock, Psychedelic, Guitar Virtuoso, male vocals",
	"The Who":                  "Rock, Hard Rock, Theatrical",
	"Black Sabbath":            "Heavy Metal, Doom",
	"Iron Maiden":              "Heavy Metal, Epic, Theatrical",
	"Judas Priest":             "Heavy Metal, Power, Leather",
	"Motorhead":                "Heavy Metal, Rock’n’Roll, Aggressive",
	"Slayer":                   "Thrash Metal, Aggressive, Dark",
	"Ozzy Osbourne":            "Heavy Metal, Dark, Theatrical, male vocals",
	"Skrillex":                 "Dubstep, Electronic, Intense, male vocals",
	"Calvin Harris":            "EDM, Dance, Festive, male vocals",
	"Avicii":                   "EDM, Melodic, Euphoric, male vocals",
	"Arctic Monkeys":           "Indie Rock, Garage, Cool",
	"Tame Impala":              "Psychedelic Rock, Dreamy, Mellifluous",
	"The Strokes":              "Indie Rock, Cool, Raw",
	"Vampire Weekend":          "Indie Rock, Eclectic, Upbeat",
	"Kings of Leon":            "Rock, Emotional, Raw",
	"The Killers":              "Rock, Synthpop, Anthemic, male vocals",
	"System of a Down":         "Metal, Political, Eccentric",
	"Radiohead":                "Alternative Rock, Experimental, Atmospheric",
	"Foo Fighters":             "Rock, Alternative, Energetic",
	"Muse":                     "Rock, Progressive, Theatrical",
	"Tool":                     "Progressive Metal, Dark, Complex",
	"Rage Against the Machine": "Rap Metal, Political, Aggressive",
	"Pearl Jam":                "Grunge, Rock, Emotional",
	"Soundgarden":              "90s Grunge, Heavy, Dark",
	"Alice in Chains":          "Grunge, Dark, Melodic",
	"Sigur Rós":                "Post-Rock, Ethereal, Atmospheric, Icelandic",
	"Björk":                    "Alternative, Experimental, Unusual, female vocals",
	"Enya":                     "New Age, Ethereal, Calm, female vocals",
	"Deadmau5":                 "Electronic, House, Progressive",
	"Marshmello":               "EDM, Dance, Happy",
	"Zedd":                     "EDM, Dance Pop, Energetic, male vocals",
	"The XX":                   "Indie Pop, Minimal, Atmospheric",
	"Lana Del Rey":             "Pop, Sadcore, Cinematic, female vocals",
	"Kacey Musgraves":          "Country, Pop, Mellifluous, female vocals",
	"St. Vincent":              "Art Rock, Eclectic, Unusual, female vocals",
	"Childish Gambino":         "HipHop, Funk, Thoughtful, male vocals",
	"SZA":                      "RnB, Neo Soul, Emotional, female vocals",
	"Frank Ocean":              "RnB, Soulful, Introspective, male vocals",
	"Tyler, The Creator":       "HipHop, Eclectic, Unusual, male vocals",
	"Solange":                  "RnB, Soul, Artistic, female vocals",
	"Brockhampton":             "HipHop, Alternative, Collective",
	"Janelle Monáe":            "Funk, RnB, Sci-Fi, female vocals",
	"Mac DeMarco":              "Indie Pop, Slacker Rock, Chill, male vocals",
	"Rufus Du Sol":             "Electronic, Dance, Atmospheric",
	"Bon Iver":                 "Indie Folk, Ethereal, Intimate, male vocals",
	"Florence + The Machine":   "Indie Rock, Dramatic, Ethereal",
	"Jack White":               "Rock, Blues, Raw, male vocals",
	"Gary Clark Jr.":           "Blues Rock, Soulful, Gritty, male vocals",
	"Leon Bridges":             "Soul, RnB, Retro, male vocals",
	"Brittany Howard":          "Rock, Soul, Powerful, female vocals",
	"Alabama Shakes":           "Rock, Blues Rock, Soulful",
	"Glass Animals":            "Psychedelic Pop, Groovy, Eclectic",
	"Portugal, The Man":        "Alternative Rock, Psychedelic, Catchy",
	"FKA Twigs":                "RnB, Electronic, Avant-Garde, female vocals",
	"The National":             "Indie Rock, Melancholy, Introspective",
	"MGMT":                     "Psychedelic Pop, Electronic, Playful",
	"Empire of the Sun":        "Electronic, Pop, Theatrical",
	"Grimes":                   "Art Pop, Electronic, Experimental, female vocals",
	"James Blake":              "Electronic, Soul, Minimalist, male vocals",
	"The War on Drugs":         "Indie Rock, Heartland Rock, Melodic",
	"Sufjan Stevens":           "Indie Folk, Baroque Pop, Intimate, male vocals",
	"Bonobo":                   "Downtempo, Electronic, Ambient",
	"Caribou":                  "Electronic, Psychedelic, Dance",
	"Four Tet":                 "Electronic, Ambient, Textural",
	"Jamie xx":                 "Electronic, House, Minimal",
	"Nicolas Jaar":             "Electronic, Experimental, Atmospheric, male vocals",
	"Flying Lotus":             "Electronic, Experimental HipHop, Fusion, male vocals",
	"Thundercat":               "Funk, Jazz, Experimental, male vocals",
	"Kamasi Washington":        "Jazz, Fusion, Epic, male vocals",
	"Massive Attack":           "Trip Hop, Dark, Atmospheric",
	"Portishead":               "Trip Hop, Dark, Cinematic",
	"Aphex Twin":               "IDM, Electronic, Experimental, male vocals",
	"Boards of Canada":         "IDM, Downtempo, Nostalgic",
	"Burial":                   "Dubstep, Ambient, Mysterious",
	"J Dilla":                  "HipHop, Soulful, Experimental, male vocals",
	"MF DOOM":                  "HipHop, Abstract, Lyrical, male vocals",
	"Blink-182":                "emo pop rock, male vocals",
	"Phoebe Bridgers":          "Bedroom, grungegaze, catchy, psychedelic, acoustic tape recording, female vocals",
}

func GenerateSongFromLyrics(prompt string, title string) ([]*model.SongResponse, string, error) {
	log.Printf("Generating song from lyrics with title: %s\n", title)

	rand.Seed(time.Now().UnixNano())
	artist := getRandomArtist()
	log.Printf("Selected artist: %s\n", artist)
	artistStyle := artists[artist]

	requestBody := map[string]interface{}{
		"is_custom":         true,
		"make_instrumental": false,
		"model_version":     "chirp-v3-5",
		"prompt":            prompt,
		"tags":              artist,
		"title":             title,
		"wait_audio":        true,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("Error marshaling request body:", err)
		return nil, "", err
	}

	sunoAPI := os.Getenv("SUNO_API") // Ensure you set this environment variable

	resp, err := http.Post(sunoAPI, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error making POST request:", err)
		return nil, "", err
	}
	defer resp.Body.Close()

	var response []model.SongResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Println("Error decoding response:", err)
		return nil, "", err
	}

	if len(response) >= 2 {
		log.Println("Successfully generated song responses.")
		return []*model.SongResponse{&response[0], &response[1]}, artistStyle, nil
	}

	log.Println("No song responses generated.")

	return nil, "", nil
}

func getRandomArtist() string {
	keys := make([]string, 0, len(artists))
	for k := range artists {
		keys = append(keys, k)
	}
	return artists[keys[rand.Intn(len(keys))]]
}
