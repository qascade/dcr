# Use official Postgres image as base
FROM postgres:latest

# Set environment variables
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD password
ENV POSTGRES_DB advertiser_db

# Copy CSV files into container
COPY ad_conversions.csv /tmp/
COPY ad_customers.csv /tmp/

# Copy initialization script to create tables and import data
COPY init.sql /docker-entrypoint-initdb.d/

# Expose default PostgreSQL port
EXPOSE 5432

# Start Postgres service
CMD ["postgres"]

