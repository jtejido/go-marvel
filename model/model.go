package model

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
	Offset  int32       `json:"offset"`
	Limit   int32       `json:"limit"`
	Total   int32       `json:"total"`
	Count   int32       `json:"count"`
	Results []Character `json:"results"`
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

type ComicDataWrapper struct {
	Code            int32              `json:"code"`
	Status          string             `json:"status"`
	Copyright       string             `json:"copyright"`
	AttributionText string             `json:"attributionText"`
	AttributionHTML string             `json:"attributionHTML"`
	Data            ComicDataContainer `json:"data"`
	Etag            string             `json:"etag"`
}

type ComicDataContainer struct {
	Offset  int32   `json:"offset"`
	Limit   int32   `json:"limit"`
	Total   int32   `json:"total"`
	Count   int32   `json:"count"`
	Results []Comic `json:"results"`
}

type Comic struct {
	ID                 int32          `json:"id"`
	DigitalID          int32          `json:"digitalId"`
	Title              string         `json:"title"`
	IssueNumber        float64        `json:"issueNumber"`
	VariantDescription string         `json:"variantDescription"`
	Description        string         `json:"description"`
	Modified           string         `json:"modified"`
	ISBN               string         `json:"isbn"`
	UPC                string         `json:"upc"`
	DiamondCode        string         `json:"diamondCode"`
	EAN                string         `json:"ean"`
	ISSN               string         `json:"issn"`
	Format             string         `json:"format"`
	PageCount          int32          `json:"pageCount"`
	TextObjects        []TextObject   `json:"textObjects"`
	ResourceURI        string         `json:"resourceURI"`
	URLs               []Url          `json:"urls"`
	Series             SeriesSummary  `json:"series"`
	Variants           []ComicSummary `json:"variants"`
	Collections        []ComicSummary `json:"collections"`
	CollectedIssues    []ComicSummary `json:"collectedIssues"`
	Dates              []ComicDate    `json:"dates"`
	Prices             []ComicPrice   `json:"prices"`
	Thumbnail          Image          `json:"thumbnail"`
	Images             []Image        `json:"images"`
	Creators           CreatorList    `json:"creators"`
	Characters         CharacterList  `json:"characters"`
	Stories            StoryList      `json:"stories"`
	Events             EventList      `json:"events"`
}

type ComicPrice struct {
	Type  string  `json:"type"`
	Price float32 `json:"price"`
}

type ComicDate struct {
	Type string `json:"type"`
	Date string `json:"data"`
}

type TextObject struct {
	Type     string `json:"type"`
	Language string `json:"language"`
	Text     string `json:"text"`
}

type CreatorList struct {
	Available     int32            `json:"available"`
	Returned      int32            `json:"returned"`
	CollectionURI string           `json:"collectionURI"`
	Items         []CreatorSummary `json:"items"`
}

type CreatorSummary struct {
	ResourceURI string `json:"resourceURI"`
	Name        string `json:"name"`
	Role        string `json:"role"`
}

type CharacterList struct {
	Available     int32              `json:"available"`
	Returned      int32              `json:"returned"`
	CollectionURI string             `json:"collectionURI"`
	Items         []CharacterSummary `json:"items"`
}

type CharacterSummary struct {
	ResourceURI string `json:"resourceURI"`
	Name        string `json:"name"`
	Role        string `json:"role"`
}
