# Gator
Gator is a linux only cli app for subscribing to rss feeds and periodically fetching posts


## Instalation
In order to install and use gator you'll first need to make sure [go](https://go.dev/doc/install) and [Postgres](https://www.postgresql.org/download/) are installed before running the command
```
go install github.com/Andrew-The-Cat/gator/cmd/gator@latest
```
In order for the app to function properly you'll also need to create a postgresql database with the name gator (you may want to use the up migrations from my github)

To finish things off you'll need to create a .gatorconfig file in your home directory with the following content
```
{
    "db_url":"postgres://[psql_username]:[psql_password]@localhost:5432/gator?sslmode=disable",
    "current_user_name":""
}
```
where the db_url will be the same as the psql url

## Usage
To use the app you'll first need to create a user with the command
```
gator register [username]
```
---
You may create multiple users and switch between them using the command
```
gator login [username]
```
---
To view a list of all users you can type
```
gator users
```
---
In order to add a feed
```
gator addfeed [feed_name] [url to rss feed endpoint]
```
Adding a feed will automatically follow said feed

---
And to view feeds
```
gator feeds
```
---
In order to grab posts from feeds other users have created you may use
```
gator follow [url]
```
---
To view follwed feeds you can run 
```
gator following
```
---
To unfollow a feed you can run
```
gator unfollow [url]
```
---
To begin fetching posts from all followed feeds 
```
gator agg [duration of time in the format XXXs/m/h]
```
The duration will determine how much the app waits between fetches
This command is meant to run in the background as gator is used within another terminal

---
To view fetched posts
```
gator browse (limit)
```
where limit is the number of posts you want to view in order from latest to oldest, default is 2