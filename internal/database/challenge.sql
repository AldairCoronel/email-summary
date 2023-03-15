-- create the transactions table
DROP TABLE IF EXISTS transactions;

CREATE TABLE transactions (
    id INTEGER PRIMARY KEY,
    date TIMESTAMP NOT NULL,
    amount FLOAT NOT NULL,
    is_credit BOOLEAN NOT NULL
);

-- create the summary table
DROP TABLE IF EXISTS summary CASCADE;

CREATE TABLE summary (
    id SERIAL PRIMARY KEY,
    total_balance FLOAT NOT NULL,
    num_of_credit_tansactions INTEGER NOT NULL,
    num_of_debit_tansactions INTEGER NOT NULL,
    total_average_credit FLOAT NOT NULL,
    total_average_debit FLOAT NOT NULL
);


-- create the month_summary table
DROP TABLE IF EXISTS month_summary;

CREATE TABLE month_summary (
    id SERIAL PRIMARY KEY,
    month VARCHAR(10) NOT NULL,
    num_of_credit_tansactions INTEGER NOT NULL,
    num_of_debit_tansactions INTEGER NOT NULL,
    average_credit FLOAT NOT NULL,
    average_debit FLOAT NOT NULL,
    summary_id SERIAL NOT NULL,
    FOREIGN KEY (summary_id) REFERENCES summary(id)
);
