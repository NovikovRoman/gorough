package main

import "log"

func main() {
	var err error

	if err = hodgepodge(); err != nil {
		log.Fatalln(err)
	}

	if err = lines(); err != nil {
		log.Fatalln(err)
	}

	if err = poligons(); err != nil {
		log.Fatalln(err)
	}

	if err = ellipses(); err != nil {
		log.Fatalln(err)
	}

	if err = curves(); err != nil {
		log.Fatalln(err)
	}

	if err = rectangles(); err != nil {
		log.Fatalln(err)
	}

	if err = arcs(); err != nil {
		log.Fatalln(err)
	}

	if err = linearPaths(); err != nil {
		log.Fatalln(err)
	}

	if err = paths(); err != nil {
		log.Fatalln(err)
	}

	if err = filling(); err != nil {
		log.Fatalln(err)
	}
}
