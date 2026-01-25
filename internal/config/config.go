package config

type Distribution struct {
	Word        int `json:"word"`
	Number      int `json:"number"`
	Period      int `json:"period"`
	Comma       int `json:"comma"`
	Quotation   int `json:"quotation"`
	Question    int `json:"question"`
	Exclamation int `json:"exclamation"`
	Brackets    int `json:"brackets"`
	Braces      int `json:"braces"`
	Parenthesis int `json:"parenthesis"`
	Colon       int `json:"colon"`
	Semicolon   int `json:"semicolon"`
}

type Config struct {
	Dictionary   string       `json:"dict"`
	TopWords     int          `json:"top"`
	WordCount    int          `json:"count"`
	Width        int          `json:"width"`
	Numbers      bool         `json:"nums"`
	Punctuation  bool         `json:"punct"`
	Distribution Distribution `json:"freqs"`
}
