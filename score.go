package p

// Score represents the Score collection in Firestore
type Score struct {
	Uid     string
	Name    string
	Points  int
	Details string
	Country string
}

// NewFromFirestoreValue returns a new Score, from values in FirestoreValue
func NewFromFirestoreValue(v FirestoreValue) (*Score, error) {
	var score Score
	uid, err := v.getStringValue("uid")
	if err != nil {
		return nil, err
	}
	name, err := v.getStringValue("name")
	if err != nil {
		return nil, err
	}
	points, err := v.getIntegerValue("points")
	if err != nil {
		return nil, err
	}
	details, err := v.getStringValue("details")
	if err != nil {
		return nil, err
	}
	country, err := v.getStringValue("country")
	if err != nil {
		return nil, err
	}
	score = Score{
		Uid:     uid,
		Name:    name,
		Points:  points,
		Details: details,
		Country: country,
	}
	return &score, nil
}
