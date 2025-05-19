package hx

type target struct {
	Val string
	Sel string
	Typ string
}

var (
	// Bodiody is the target element for the body of the page.
	//
	// The body of the page is the whole page excluding the header.
	Bodiody = target{
		Val: "bodiody",
		Sel: "#bodiody",
		Typ: "id",
	}

	// Taggle is the target element for the content of the tag page.
	Taggle = target{
		Val: "taggle",
		Sel: "#taggle",
		Typ: "id",
	}
)
