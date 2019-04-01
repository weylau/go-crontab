package main

import "mcrontab/master/app"

func main() {
	var (
		master *app.App
	)
	master = app.NewApp()
	master.Run()
}
