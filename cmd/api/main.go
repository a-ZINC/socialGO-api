package main

func main() {
	app := &Application{
		Config: Config{
			Addr: ":4000",
		},
	}
	mux := app.mount()

	if err := app.run(mux); err != nil {
		panic(err)
	}
}