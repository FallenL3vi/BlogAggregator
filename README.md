## REQUIREMENTS
- PostgreSQL 14+
- Go 1.22+


## Instalation instruction
- Run this in the terminal
go install https://github.com/FallenL3vi/BlogAggregator

Create file named .gatorconfig.json in your user home direction
Write inside:
"db_url":"postgres://postgres:postgres@localhost:5432/blogaggre?sslmode=disable"
"current_user_name":"CURRENT_USER_NAME"

To run the program go inside program direction
Use command:
go run . [command]

Commands:
reset -> Resets tables
register [user name] -> to register and swap to new user
addfeed [Feed title] [URL] -> add feed to list of feeds adn follow it
login [user name] -> to change user
following -> prints list of followed feeds by current user
follow [URL] -> to follow feed
agg [1s/m/h] -> start aumotaic aggregation from feeds followed by current user. It does it every time stamp included by user
browse [LIMIT] -> Prints Title and Description of articles from feeds. Limited by default to 2



