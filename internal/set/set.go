package set

type Strings map[string]struct{}

func OfStrings(ss []string) Strings {
	set := make(Strings)
	for _, s := range ss {
		set[s] = struct{}{}
	}
	return set
}
