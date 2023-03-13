DROP TABLE IF EXISTS transactions;

CREATE TABLE transactions (
    id INTEGER PRIMARY KEY,
    date TIMESTAMP NOT NULL,
    amount FLOAT NOT NULL,
    is_credit BOOLEAN NOT NULL,
);

DROP TABLE IF EXISTS month_summary;

CREATE TABLE month_summary (
    id SERIAL PRIMARY KEY,
    year INTEGER,
    month VARCHAR(10) NOT NULL,
    numOfTransactions INTEGER NOT NULL,
    averageCredit FLOAT NOT NULL,
    averageDebit FLOAT NOT NULL
);

DROP TABLE IF EXISTS summary;

CREATE TABLE summary (
    id SERIAL PRIMARY KEY,
    totalBalance FLOAT NOT NULL,
    numOfTotalTransactions INTEGER NOT NULL,
    monthlySummaries JSONB NOT NULL
);
