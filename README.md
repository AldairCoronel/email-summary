# Stori Challenge

This project is a solution for the Stori Challenge. It is written in Go and uses a PostgreSQL database.

## Getting Started
To run this project, you will need to have Docker installed. Follow the instructions for your operating system on the Docker website to install Docker.

1. First, pull the Docker image from Docker Hub:

```
docker pull aldacacr/stori-challenge:1.0.0
```

2. Run the Docker container:

```
docker run -it aldacacr/stori-challenge:1.0.0
```


You will access the command line and I need you to do the last database configuration steps:

3. Execute the following in ``root@id:/app#``:

```
su postgres
```

4. Execute Postgres in ``postgres@id:/app$``:

```
psql
```

5. Create the challenge database:

```
CREATE DATABASE challenge;
```

6. Connect to the challenge database:

```
\c challenge
```

7, Create the tables from the .sql file:

```
\i internal/database/challenge.sql
```

8. Set the Postgres password:

```
alter user postgres password 'hola';
```

9. Exit the Postgres prompt and the container:
```
exit
exit
```

10. Now you can execute main.go passing a flag with your email in this path ``root@id:/app#``:

```
go run cmd/main.go -- emailTo <your.email@example.com>
```

This should give you

```
New account created with ID: 1
Email summary sent!
```

The email has been sent.


## License
This project is licensed under the MIT License - see the LICENSE file for details.
