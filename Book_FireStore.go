package main

type Book_FireStore struct {
	Id			int64	 `firestore:"id"`
	BookName   	string   `firestore:"bookName,omitempty"`
	Genre      	string   `firestore:"genre,omitempty"`
	SubGenre   	string   `firestore:"subGenre,omitempty"`
	ISBN    	string   `firestore:"ISBN,omitempty"`
	Find 		int64    `firestore:"find,omitempty"`
	Sum 		int64    `firestore:"sum,omitempty"`
	Author    	string   `firestore:"author,omitempty"`
	Publisher   string   `firestore:"publisher,omitempty"`
	Pubdate		string   `firestore:"pubdate,omitempty"`
	Exist 		string   `firestore:"exist,omitempty"`
	LocateAt4F  bool     `firestore:"locateAt4F,omitempty"`
	WithDisc    string   `firestore:"withDisc,omitempty"`
	Other    	string   `firestore:"other,omitempty"`
	Borrower	interface{} `firestore:"borrower,omitempty"`
	Location	string 	 `firestore:"location,omitempty"`
}
