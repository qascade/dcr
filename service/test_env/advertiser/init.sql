-- Create ad_conversions table

CREATE SCHEMA advertiser_db;
-- revoke usage on schema public from public;
SET search_path TO  'advertiser_db';
-- SET SCHEMA 'advertiser_db';

CREATE TABLE advertiser_db.ad_conversions (
    email varchar(255),
    product varchar(255),
    sls_date date,
    sales_dlr numeric
);

-- Import data from ad_conversions.csv
COPY ad_conversions FROM '/tmp/ad_conversions.csv' DELIMITER ',' CSV HEADER;

-- Create ad_customers table
CREATE TABLE ad_customers (
    email varchar(255),
    phone varchar(20),
    pets varchar(255),
    zip varchar(10)
);

-- Import data from ad_customers.csv
COPY ad_customers FROM '/tmp/ad_customers.csv' DELIMITER ',' CSV HEADER;
