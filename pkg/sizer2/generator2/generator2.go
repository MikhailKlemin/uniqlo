package generator2

//Opts for Generate
type Opts struct {
	Gender     string //Men
	Weight     int
	Height     int
	Shape      int
	Chest      int
	Age        int
	Preference int
}

//Generate is
func Generate(o Opts) []Opts {
	var out []Opts

	makeRange := func(min, max int) []int {
		a := make([]int, max-min+1)
		for i := range a {
			a[i] = min + i
		}
		return a
	}
	sm := make(map[int][]int)
	sm[140] = makeRange(40, 65)
	sm[145] = makeRange(43, 65)
	sm[155] = makeRange(40, 120)
	sm[160] = makeRange(44, 121)
	sm[165] = makeRange(45, 123)
	sm[165] = makeRange(45, 123)
	sm[165] = makeRange(45, 123)
	sm[165] = makeRange(45, 123)
	sm[170] = makeRange(50, 125)
	sm[175] = makeRange(50, 128)
	sm[180] = makeRange(50, 130)
	sm[185] = makeRange(55, 130)
	sm[190] = makeRange(55, 130)
	sm[195] = makeRange(60, 135)
	sm[200] = makeRange(65, 135)
	sm[205] = makeRange(70, 132)
	sm[210] = makeRange(80, 132)

	for height, weights := range sm {
		for _, weight := range weights {
			var c Opts
			c.Gender = o.Gender
			c.Weight = weight
			c.Height = height
			c.Age = o.Age
			c.Chest = o.Chest
			c.Preference = o.Preference
			c.Shape = o.Shape
			out = append(out, c)
		}
	}

	return out
}
