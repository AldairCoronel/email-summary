-- create the accounts table
DROP TABLE IF EXISTS accounts CASCADE;

CREATE TABLE accounts (
    account_id SERIAL PRIMARY KEY
);

-- create the transactions table
DROP TABLE IF EXISTS transactions;

CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id SERIAL NOT NULL,
    id INTEGER NOT NULL,
    date TIMESTAMP NOT NULL,
    amount FLOAT NOT NULL,
    is_credit BOOLEAN NOT NULL,
    FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

-- create the summary table
DROP TABLE IF EXISTS summary CASCADE;

CREATE TABLE summary (
    summary_id SERIAL PRIMARY KEY,
    account_id SERIAL NOT NULL,
    total_balance FLOAT NOT NULL,
    total_transactions INTEGER NOT NULL,
    num_of_credit_transactions INTEGER NOT NULL,
    num_of_debit_transactions INTEGER NOT NULL,
    total_average_credit FLOAT NOT NULL,
    total_average_debit FLOAT NOT NULL,
    FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

-- create the month_summary table
DROP TABLE IF EXISTS month_summary;

CREATE TABLE month_summary (
    month_summary_id SERIAL PRIMARY KEY,
    month VARCHAR(10) NOT NULL,
    total_balance FLOAT NOT NULL,
    total_transactions INTEGER NOT NULL,
    num_of_credit_transactions INTEGER NOT NULL,
    num_of_debit_transactions INTEGER NOT NULL,
    average_credit FLOAT NOT NULL,
    average_debit FLOAT NOT NULL,
    summary_id SERIAL NOT NULL,
    FOREIGN KEY (summary_id) REFERENCES summary(summary_id)
);
