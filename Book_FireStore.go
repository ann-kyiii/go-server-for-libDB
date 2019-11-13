package main

type BookFireStore struct {
	Id			int64	 `firestore:"id"`
	BookName   	string   `firestore:"bookName"`
	Genre      	string   `firestore:"genre"`
	SubGenre   	string   `firestore:"subGenre"`
	ISBN    	string   `firestore:"ISBN"`
	Find 		int64    `firestore:"find"`
	Sum 		int64    `firestore:"sum"`
	Author    	string   `firestore:"author"`
	Publisher   string   `firestore:"publisher"`
	Pubdate		string   `firestore:"pubdate"`
	Exist 		string   `firestore:"exist"`
	LocateAt4F  bool     `firestore:"locateAt4F"`
	WithDisc    string   `firestore:"withDisc"`
	Other    	string   `firestore:"other"`
	Borrower	interface{} `firestore:"borrower"`
	Location	string 	 `firestore:"location"`
	ImgURL		string	 `firestore:"imgURL"`	
}
