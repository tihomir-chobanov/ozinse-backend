package model

type Episode struct {
	ID             int    `json:"id"`
	SeasonID       int    `json:"season_id"`
	EpisodeNumber  int    `json:"episode_number"`
	YoutubeVideoID string `json:"youtube_video_id"`
	Duration       int    `json:"duration"`
}

/*

The Model is the simplest but most essential layer. It defines exactly what your data looks like as it travels through the Handlers, Services, and Repositories.

Responsibilities:
Data Structure: It defines the properties of an entity (e.g., a Category has an ID and a Name).

Tagging: It uses "tags" (like `json:"name"`) to tell Go how to rename or format fields when converting them to JSON for the client.

Consistency: By using the same Model across all layers, you ensure that the "Chef" (Service) and the "Supplier" (Repository) are always talking about the same "Menu Item."

Why it's important: It serves as the "single source of truth." If you add a new field to your database (like description), you update the Model once, and every other layer immediately knows how to handle that new data.

*/
