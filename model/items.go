package model

type ModelItem interface{
	Print() string
	Equals(m ModelItem) bool
}

type User struct{
	Id float64
	FirstName, LastName, Username, Email, Password string
}

type CloudHandler interface{
	Init() string
	GetFileList() map[int]string
	DownloadById(id int)
	DownloadByName(name string)
	Upload(path, name string)
}

func (u *User) Print() string{
	return u.FirstName + ", " + u.LastName + ", " + u.Username + ", " + u.Email + ", " + u.Password
}

func (u *User) Equals(m ModelItem) bool{
	return u.Print() == m.Print()
}



