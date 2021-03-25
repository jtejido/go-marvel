package model

type Response struct {
	Etag   string               `json:"etag"`
	Result CharacterDataWrapper `json:"result"`
}

type CharacterDataWrapper struct {
	Code            int32                  `json:"code"`
	Status          string                 `json:"status"`
	Copyright       string                 `json:"copyright"`
	AttributionText string                 `json:"attributionText"`
	AttributionHTML string                 `json:"attributionHTML"`
	Data            CharacterDataContainer `json:"data"`
	Etag            string                 `json:"etag"`
}

type CharacterDataContainer struct {
	Offset  int32        `json:"offset"`
	Limit   int32        `json:"limit"`
	Total   int32        `json:"total"`
	Count   int32        `json:"count"`
	Results []*Character `json:"results"`
}

type Character struct {
	ID          int32      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Modified    string     `json:"modified"`
	ResourceURI string     `json:"resourceURI"`
	Thumbnail   Image      `json:"thumbnail"`
	URLs        []Url      `json:"urls"`
	Comics      ComicList  `json:"comics"`
	Stories     StoryList  `json:"stories"`
	Events      EventList  `json:"events"`
	Series      SeriesList `json:"series"`
}

type Image struct {
	Path      string `json:"path"`
	Extension string `json:"extension"`
}

type Url struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type ComicList struct {
	Available     int32          `json:"available"`
	Returned      int32          `json:"returned"`
	CollectionURI string         `json:"collectionURI"`
	Items         []ComicSummary `json:"items"`
}

type ComicSummary struct {
	ResourceURI string `json:"resourceURI"`
	Name        string `json:"name"`
}

type StoryList struct {
	Available     int32          `json:"available"`
	Returned      int32          `json:"returned"`
	CollectionURI string         `json:"collectionURI"`
	Items         []StorySummary `json:"items"`
}

type StorySummary struct {
	ResourceURI string `json:"resourceURI"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

type EventList struct {
	Available     int32          `json:"available"`
	Returned      int32          `json:"returned"`
	CollectionURI string         `json:"collectionURI"`
	Items         []EventSummary `json:"items"`
}

type EventSummary struct {
	ResourceURI string `json:"resourceURI"`
	Name        string `json:"name"`
}

type SeriesList struct {
	Available     int32           `json:"available"`
	Returned      int32           `json:"returned"`
	CollectionURI string          `json:"collectionURI"`
	Items         []SeriesSummary `json:"items"`
}

type SeriesSummary struct {
	ResourceURI string `json:"resourceURI"`
	Name        string `json:"name"`
}
