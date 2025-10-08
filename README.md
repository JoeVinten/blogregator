## Blogregator

### What this project is about
A simple command line blog aggregator written in Go. It follows one the the guided courses on [boot.dev](https://boot.dev) to learn Go and backend development as a whole. While the vast majority of the code is my own I have gone back over the solutions provided by boot.dev to learn from them and make sure I understand everything and refactored inline with this.

### How to run this project
You'll need to have Postgres and Go installed to run this project as it's currently only runs locally and uses a Postgres database.

You can install the blogregator project by cloning this repository, and then running `go build` to get the binary. After this run `go install`.

You should then be able to access the commands by running `blogregator <command>`.


### Config

You'll need a `.gatorconfig.json` file in your home directory with the following json structure:

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the values with your database connection string.


### Running commands

Then you can run the CLI commands by using `blogregator <command>`.



