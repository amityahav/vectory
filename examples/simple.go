package examples

import (
	"Vectory/db"
	"context"
)

func main() {
	// just checking the api for undesired exported fields
	d, _ := db.Open("")
	c, _ := d.CreateCollection(context.TODO(), nil)
	print(c)
}
