package membership

type UserLister interface {
	At(int) User
	Len() int
}
