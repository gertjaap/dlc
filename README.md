# DLC Oracle

This project can serve as an oracle while forming Discreet Log Contracts. This oracle currently publishes the value of the US Dollar denominated in Bitcoin's smallest fraction (satoshis). You can interact with the oracle via simple REST calls. A live version of this oracle is running on [https://oracle.gertjaap.org/] 

If you want to learn more about Discreet Log Contracts, checkout the [whitepaper](https://adiabat.github.io/dlc.pdf)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

You need to have golang installed, or you can use Docker

### Installing

First, clone the repository and install the dependencies

```
git clone https://github.com/gertjaap/dlcoracle
cd dlcoracle
go get -v ./...
```

Then you can build the oracle using
```
go build
```

### PostgreSQL backend 

The datasource(s) that the oracle offers keys and signatures for are stored in a PostgreSQL database. You'll have to set the environment variable `DLC_DB_CONN_STRING` to a valid
PostgreSQL server connection string, by typing the following in your shell: `export DLC_DB_CONN_STRING=postgres://postgres:password@localhost/database?sslmode=disable`. Note that this is just an example connection string - be sure to update the username, password, hostname and database to match your database server!

Once the connection string is set, you'll need to create a table in the PostgreSQL database you've specified, and  name it `datasources`. The table will need to have, at a minimum, the following fields:

```
                Table "public.datasources"
   Column    |            Type             |   Modifiers   
-------------+-----------------------------+---------------
 id          | integer                     | not null
 name        | text                        | not null
 description | text                        | not null
 value       | integer                     | not null
 interval    | integer                     | not null
```

To create the table, connect to your database and enter the following command:
```
psql -h localhost -U postgres -d database
psql (9.6.7)
Type "help" for help.
database=# create table datasources ( Id int primary key not null, Name text not null, Description text not null, Value int not null, Interval int not null);
```

You'll also need to populate the table with at least one test datasource:
```
database=# insert into datasources values (1, 'US Dollar', 'Publishes the value of USD denominated in 1/100000000th units of BTC (satoshi) in multitudes of 10', '14200', '300');
```


### Running the oracle

Simply start the executable. Since the oracle generates a private key it will ask you for a password to protect it, that you have to enter each time you start up the oracle.

```
./dlcoracle
```

## REST Endpoints

| resource          | description                              |
|:------------------|:-----------------------------------------|
|[`/api/pubkey`](https://oracle.gertjaap.org/api/pubkey)      | Returns the public keys of the oracle     |
|[`/api/datasources`](https://oracle.gertjaap.org/api/datasources) | Returns an array of data sources the oracle publishes |
|[`/api/rpoint/{s}/{t}`](https://oracle.gertjaap.org/api/rpoint/1/1523447385) | Returns the public one-time-signing key for datasource with ID **s** at the unix timestamp **t**. |
|[`/api/publication/{R}`](https://oracle.gertjaap.org/api/publication/1/1523447385) | Returns the value and signature published for data source point **R** (if published). R is hex encoded [33]byte |

## Using the public deployment

You're free to use my public deployment of the oracle as well. I have linked the URLs of the public deployment in the REST endpoint table above.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
