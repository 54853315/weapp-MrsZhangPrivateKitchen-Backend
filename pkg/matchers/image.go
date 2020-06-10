package matchers

//type MatcherImages struct{
//	Name string
//	Ext string
//	MIME string
//}

var imageMime = []struct {
	Name string
	Ext  string
	//MIME string
}{
	{"image/jpeg", "jpg"},
	{"image/jp2", "jp2"},
	{"image/png", "png"},
	{"image/gif", "gif"},
}

func Image(mime string) bool {
	for _, target := range imageMime {
		if target.Name == mime {
			return true
		}
	}
	return false
}
