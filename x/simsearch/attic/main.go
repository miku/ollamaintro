package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// Movie represents a movie with its metadata and embedding
type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Genre       string    `json:"genre"`
	Year        int       `json:"year"`
	Plot        string    `json:"plot"`
	Embedding   []float64 `json:"embedding,omitempty"`
	EmbeddingID string    `json:"embedding_id,omitempty"`
}

// MovieDatabase holds our collection of movies and their embeddings
type MovieDatabase struct {
	Movies     []Movie           `json:"movies"`
	Index      map[string][]int  `json:"index"` // genre -> movie indices
	Embeddings map[int][]float64 `json:"embeddings"`
}

// OllamaEmbeddingRequest represents the request to Ollama
type OllamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaEmbeddingResponse represents the response from Ollama
type OllamaEmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// SimilarityResult holds a movie and its similarity score
type SimilarityResult struct {
	Movie      Movie
	Similarity float64
}

const (
	OLLAMA_URL = "http://k9:11434"
	MODEL_NAME = "nomic-embed-text" // or "all-minilm" - adjust based on your Ollama setup
	DATA_FILE  = "movie_embeddings.gob"
)

func main() {
	db := &MovieDatabase{
		Movies:     make([]Movie, 0),
		Index:      make(map[string][]int),
		Embeddings: make(map[int][]float64),
	}

	// Load existing data or create sample dataset
	if err := db.Load(); err != nil {
		fmt.Println("Creating new movie database...")
		db.createSampleMovies()

		fmt.Println("Generating embeddings...")
		if err := db.generateEmbeddings(); err != nil {
			log.Fatal("Failed to generate embeddings:", err)
		}

		fmt.Println("Saving database...")
		if err := db.Save(); err != nil {
			log.Fatal("Failed to save database:", err)
		}
	}

	// Interactive demo
	fmt.Printf("\n🎬 Movie Similarity Search Demo\n")
	fmt.Printf("Database loaded with %d movies\n\n", len(db.Movies))

	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Find similar movies to a specific movie")
		fmt.Println("2. Search by custom plot description")
		fmt.Println("3. Show genre clusters")
		fmt.Println("4. Exit")
		fmt.Print("\nChoice: ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			db.findSimilarMovies()
		case "2":
			db.searchByDescription()
		case "3":
			db.showGenreClusters()
		case "4":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
		fmt.Println()
	}
}

func (db *MovieDatabase) createSampleMovies() {
	movies := []Movie{
		{ID: 1, Title: "The Matrix", Genre: "Sci-Fi", Year: 1999, Plot: "A computer hacker learns from mysterious rebels about the true nature of his reality and his role in the war against its controllers."},
		{ID: 2, Title: "Inception", Genre: "Sci-Fi", Year: 2010, Plot: "A thief who steals corporate secrets through dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O."},
		{ID: 3, Title: "The Terminator", Genre: "Sci-Fi", Year: 1984, Plot: "A human soldier is sent from 2029 to 1984 to stop an almost indestructible cyborg killing machine, sent from the same year, which has been programmed to execute a young woman whose unborn son is the key to humanity's future salvation."},
		{ID: 4, Title: "Blade Runner", Genre: "Sci-Fi", Year: 1982, Plot: "A blade runner must pursue and terminate four replicants who stole a ship in space, and have returned to the earth seeking their creator."},
		{ID: 5, Title: "The Shawshank Redemption", Genre: "Drama", Year: 1994, Plot: "Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency."},
		{ID: 6, Title: "Forrest Gump", Genre: "Drama", Year: 1994, Plot: "The presidencies of Kennedy and Johnson, the events of Vietnam, Watergate and other historical events unfold through the perspective of an Alabama man with an IQ of 75."},
		{ID: 7, Title: "The Godfather", Genre: "Crime", Year: 1972, Plot: "The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son."},
		{ID: 8, Title: "Goodfellas", Genre: "Crime", Year: 1990, Plot: "The story of Henry Hill and his life in the mob, covering his relationship with his wife Karen Hill and his mob partners Jimmy Conway and Tommy DeVito."},
		{ID: 9, Title: "Pulp Fiction", Genre: "Crime", Year: 1994, Plot: "The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption."},
		{ID: 10, Title: "The Dark Knight", Genre: "Action", Year: 2008, Plot: "When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice."},
		{ID: 11, Title: "Die Hard", Genre: "Action", Year: 1988, Plot: "An NYPD officer tries to save his wife and several others taken hostage by German terrorists during a Christmas party at the Nakatomi Plaza in Los Angeles."},
		{ID: 12, Title: "Mad Max: Fury Road", Genre: "Action", Year: 2015, Plot: "In a post-apocalyptic wasteland, a woman rebels against a tyrannical ruler in search for her homeland with the aid of a group of female prisoners, a psychotic worshiper, and a drifter named Max."},
		{ID: 13, Title: "The Princess Bride", Genre: "Fantasy", Year: 1987, Plot: "A bedridden boy's grandfather reads him the story of a farmboy-turned-pirate who encounters numerous obstacles, enemies and allies in his quest to be reunited with his true love."},
		{ID: 14, Title: "Lord of the Rings: Fellowship", Genre: "Fantasy", Year: 2001, Plot: "A meek Hobbit from the Shire and eight companions set out on a journey to destroy the powerful One Ring and save Middle-earth from the Dark Lord Sauron."},
		{ID: 15, Title: "Harry Potter: Sorcerer's Stone", Genre: "Fantasy", Year: 2001, Plot: "An orphaned boy enrolls in a school of wizardry, where he learns the truth about himself, his family and the terrible evil that haunts the magical world."},
		{ID: 16, Title: "Titanic", Genre: "Romance", Year: 1997, Plot: "A seventeen-year-old aristocrat falls in love with a kind but poor artist aboard the luxurious, ill-fated R.M.S. Titanic."},
		{ID: 17, Title: "Casablanca", Genre: "Romance", Year: 1942, Plot: "A cynical nightclub owner protects an old flame and her husband from Nazis in Morocco."},
		{ID: 18, Title: "The Notebook", Genre: "Romance", Year: 2004, Plot: "A poor yet passionate young man falls in love with a rich young woman, giving her a sense of freedom, but they are soon separated because of their social differences."},
		{ID: 19, Title: "Alien", Genre: "Horror", Year: 1979, Plot: "After a space merchant vessel receives an unknown transmission as a distress call, one of the crew is attacked by a mysterious life form and they soon realize that its life cycle has merely begun."},
		{ID: 20, Title: "The Exorcist", Genre: "Horror", Year: 1973, Plot: "When a teenage girl is possessed by a mysterious entity, her mother seeks the help of two priests to save her daughter."},
		{ID: 21, Title: "Get Out", Genre: "Horror", Year: 2017, Plot: "A young African-American visits his white girlfriend's parents for the weekend, where his simmering uneasiness about their reception of him eventually reaches a boiling point."},
		{ID: 22, Title: "Toy Story", Genre: "Animation", Year: 1995, Plot: "A cowboy doll is profoundly threatened and jealous when a new spaceman figure supplants him as top toy in a boy's room."},
		{ID: 23, Title: "Finding Nemo", Genre: "Animation", Year: 2003, Plot: "After his son is captured in the Great Barrier Reef and taken to Sydney, a timid clownfish sets out on a journey to bring him home."},
		{ID: 24, Title: "The Lion King", Genre: "Animation", Year: 1994, Plot: "A Lion prince is cast out of his pride by his cruel uncle, who claims he killed his father so that he can become the new king."},
		{ID: 25, Title: "Interstellar", Genre: "Sci-Fi", Year: 2014, Plot: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival."},
		{ID: 26, Title: "Gravity", Genre: "Sci-Fi", Year: 2013, Plot: "Two astronauts work together to survive after an accident leaves them stranded in space."},
		{ID: 27, Title: "2001: A Space Odyssey", Genre: "Sci-Fi", Year: 1968, Plot: "After discovering a mysterious artifact buried beneath the Lunar surface, mankind sets off on a quest to find its origins with help from intelligent supercomputer H.A.L. 9000."},
		{ID: 28, Title: "E.T.", Genre: "Sci-Fi", Year: 1982, Plot: "A troubled child summons the courage to help a friendly alien escape Earth and return to his home world."},
		{ID: 29, Title: "Casino Royale", Genre: "Action", Year: 2006, Plot: "After earning 00 status and a licence to kill, Secret Agent James Bond sets out on his first mission as 007. Bond must defeat a private banker funding terrorists in a high-stakes game of poker."},
		{ID: 30, Title: "Mission: Impossible", Genre: "Action", Year: 1996, Plot: "An American agent, under false suspicion of disloyalty, must discover and expose the real spy without the help of his organization."},
	}

	db.Movies = movies

	// Build genre index
	for i, movie := range movies {
		if _, exists := db.Index[movie.Genre]; !exists {
			db.Index[movie.Genre] = make([]int, 0)
		}
		db.Index[movie.Genre] = append(db.Index[movie.Genre], i)
	}
}

func (db *MovieDatabase) generateEmbeddings() error {
	for i := range db.Movies {
		movie := &db.Movies[i]

		// Create a rich text representation for embedding
		text := fmt.Sprintf("Title: %s. Genre: %s. Plot: %s",
			movie.Title, movie.Genre, movie.Plot)

		embedding, err := getEmbedding(text)
		if err != nil {
			return fmt.Errorf("failed to get embedding for movie %s: %v", movie.Title, err)
		}

		movie.Embedding = embedding
		db.Embeddings[movie.ID] = embedding

		fmt.Printf("Generated embedding for: %s\n", movie.Title)

		// Be nice to the API
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func getEmbedding(text string) ([]float64, error) {
	reqBody := OllamaEmbeddingRequest{
		Model:  MODEL_NAME,
		Prompt: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(OLLAMA_URL+"/api/embeddings", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response OllamaEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedding, nil
}

func (db *MovieDatabase) findSimilarMovies() {
	fmt.Println("\nAvailable movies:")
	for i, movie := range db.Movies {
		fmt.Printf("%d. %s (%s, %d)\n", i+1, movie.Title, movie.Genre, movie.Year)
	}

	fmt.Print("\nEnter movie number: ")
	var choice int
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(db.Movies) {
		fmt.Println("Invalid choice.")
		return
	}

	targetMovie := db.Movies[choice-1]
	results := db.findSimilar(targetMovie.Embedding, 5)

	fmt.Printf("\nMovies similar to '%s':\n", targetMovie.Title)
	for i, result := range results {
		if result.Movie.ID != targetMovie.ID { // Exclude the target movie itself
			fmt.Printf("%d. %s (%.3f similarity)\n   Genre: %s | %s\n",
				i+1, result.Movie.Title, result.Similarity, result.Movie.Genre, result.Movie.Plot[:100]+"...")
		}
	}
}

func (db *MovieDatabase) searchByDescription() {
	fmt.Print("\nEnter a plot description: ")
	var description string
	fmt.Scanln(&description)

	if strings.TrimSpace(description) == "" {
		fmt.Println("Please enter a description.")
		return
	}

	fmt.Println("Generating embedding for your description...")
	embedding, err := getEmbedding(description)
	if err != nil {
		fmt.Printf("Error generating embedding: %v\n", err)
		return
	}

	results := db.findSimilar(embedding, 5)

	fmt.Printf("\nMovies matching your description:\n")
	for i, result := range results {
		fmt.Printf("%d. %s (%.3f similarity)\n   Genre: %s | %s\n",
			i+1, result.Movie.Title, result.Similarity, result.Movie.Genre, result.Movie.Plot[:100]+"...")
	}
}

func (db *MovieDatabase) showGenreClusters() {
	fmt.Println("\nGenre Clusters:")

	for genre, indices := range db.Index {
		fmt.Printf("\n%s Movies:\n", genre)

		// Calculate average inter-cluster similarity
		var totalSim float64
		var count int

		for i, idx1 := range indices {
			for j, idx2 := range indices {
				if i < j {
					sim := cosineSimilarity(db.Movies[idx1].Embedding, db.Movies[idx2].Embedding)
					totalSim += sim
					count++
				}
			}
		}

		avgSim := totalSim / float64(count)
		fmt.Printf("  Average intra-genre similarity: %.3f\n", avgSim)

		for _, idx := range indices {
			movie := db.Movies[idx]
			fmt.Printf("  - %s (%d)\n", movie.Title, movie.Year)
		}
	}
}

func (db *MovieDatabase) findSimilar(targetEmbedding []float64, count int) []SimilarityResult {
	var results []SimilarityResult

	for _, movie := range db.Movies {
		if len(movie.Embedding) == 0 {
			continue
		}

		similarity := cosineSimilarity(targetEmbedding, movie.Embedding)
		results = append(results, SimilarityResult{
			Movie:      movie,
			Similarity: similarity,
		})
	}

	// Sort by similarity (highest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	if len(results) > count {
		results = results[:count]
	}

	return results
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func (db *MovieDatabase) Save() error {
	file, err := os.Create(DATA_FILE)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(db)
}

func (db *MovieDatabase) Load() error {
	file, err := os.Open(DATA_FILE)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(db)
}
