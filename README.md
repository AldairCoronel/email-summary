# Stori Challenge

This project is a solution for the Stori Challenge. It is written in Go and uses a PostgreSQL database.

## Getting Started
To run this project, you will need to have Docker installed. Follow the instructions for your operating system on the Docker website to install Docker.

1. First, pull the Docker image from Docker Hub:

```
docker pull aldacacr/stori-challenge:1.0.1
```

2. Run the Docker container:

```
docker run -it aldacacr/stori-challenge:1.0.1
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
go run cmd/main.go --csv="./sample/txns.csv" --emailTo="<your.email@example.com>"
```


This should give you

```
New account created with ID: 1
Email summary sent!
```

The email has been sent.

## Structure

```
email-summary/
├── cmd
│   └── main.go
├── Dockerfile
├── go.mod
├── go.sum
├── internal
│   ├── controller
│   │   └── controller.go
│   ├── database
│   │   ├── challenge.sql
│   │   └── postgres.go
│   ├── models
│   │   ├── account.go
│   │   ├── summary.go
│   │   └── transaction.go
│   ├── repository
│   │   └── repository.go
│   └── view
│       ├── email.go
│       └── email-template.html
├── README.md
├── sample
│   └── txns.csv
└── Technical_Challenge_Stori.pdf
```


The project is structured as follows:


`cmd`: contains the `main.go` file which serves as the entry point for the application.

`Dockerfile`: contains instructions for building a Docker image of the application.

`go.mod` and `go.sum`: files that define the dependencies of the project.

`internal`: contains the internal packages of the application.

`controller`: contains the TransactionController which is responsible for creating an account, processing a CSV file, computing the summary, and storing the info in the database.

`database`: contains the PostgresRepository which implements the Repository interface for storing and retrieving data from the PostgreSQL database. And the file for the tables creation.

`models`: contains the Account, Summary, and Transaction models which define the structures of the data used in the application.

`repository`: contains the Repository interface which defines the methods for storing and retrieving data from the database.

`view`: contains the SMTPService which implements the EmailService interface for sending email summaries to the specified email address.

`sample`: contains an example CSV file.

`Technical_Challenge_Stori.pdf`: contains the challenge requirements.


The application follows a repository design pattern where the controllers and views are independent of the data storage implementation. The TransactionController interacts with the Repository interface to store and retrieve data from the database, while the SMTPService interacts with the EmailService interface to send email summaries. The PostgresRepository implements the Repository interface using the PostgreSQL database. The SMTPService uses Gmail as the email service and a pre-defined HTML template to style the email summary. The main.go file serves as the entry point for the application, which initializes the required dependencies and calls the TransactionController to create an account, process a CSV file, compute the summary, and send an email summary.


## License
This project is licensed under the MIT License - see the LICENSE file for details.
