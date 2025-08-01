## Helpful Commands

**docker compose exec -it db psql -U baloo -d lenslocked**

```
docker compose exec is like docker exec, which runs a command within
a docker container. Key difference is that docker **compose** exec also 
allows for referring to the service listed in the compose.yml file
```

Note the code below

```yaml
version: "3.9"
services:
  # Our Postgres database
  db:
    # The service will be named db.
    image: postgres # The postgres image will be used
    restart: always # Always try to restart if this stops running
    environment:
      # Provide environment variables
      POSTGRES_USER: baloo # POSTGRES_USER env var w/ value baloo
      POSTGRES_PASSWORD: junglebook
      POSTGRES_DB: lenslocked
    ports:
      # Expose ports so that apps not running via docker compose can connect to them.
      - 5432:5432 # format here is "port on our machine":"port on container"
  # Adminer provides a nice little web UI to connect to databases
  adminer:
    image: adminer
    restart: always
    environment:
      ADMINER_DESIGN: dracula # Pick a theme - https://github.com/vrana/adminer/tree/master/desi
    ports:
      - 3333:8080
```
We are referring to the **db** service that our config is firing up

*Additional Flags Used*
* -i - interactive, means that the container can be dealt with from **WITHIN** the container
* -t - allows for a nice terminal setup

so in plain english the command is saying: 
```
* execute the psql command
* within the db service or container
* interactively, using a terminal
* connecting as user baloo, to the database lenslocked
```