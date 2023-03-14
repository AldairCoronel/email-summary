DROP TABLE IF EXISTS transactions;

CREATE TABLE transactions (
    id INTEGER PRIMARY KEY,
    date TIMESTAMP NOT NULL,
    amount FLOAT NOT NULL,
    is_credit BOOLEAN NOT NULL
);

DROP TABLE IF EXISTS month_summary;

CREATE TABLE month_summary (
    month VARCHAR(10) PRIMARY KEY,
    year INTEGER NOT NULL,
    num_of_transactions INTEGER NOT NULL,
    average_credit FLOAT NOT NULL,
    average_debit FLOAT NOT NULL,
    summary_id  INTEGER NOT NULL REFERENCES summary(id)
);

DROP TABLE IF EXISTS summary CASCADE;

CREATE TABLE summary (
    id SERIAL PRIMARY KEY,
    total_balance FLOAT NOT NULL,
    num_of_total_transactions INTEGER NOT NULL,
    monthly_summaries JSONB NOT NULL
);
