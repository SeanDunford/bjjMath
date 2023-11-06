package scraper

type Circle struct {
	x int
	y int
}

type Rectangle struct {
	x int
	y int
}

type Shape interface {
	Circle | Rectangle
	area()
}
