# common customer count 
select count(distinct email) from media_customers as m join ad_customers as a on a.email = m.email

select count({{}}) from {{table1}} as {{}} join {{table2}} as {{}} on {{}} = {{}}



# advertiser 
select media_exposures.campaign, count(*) as impressions, count(distinct ad_conversions.email) as conversions from media_exposures left join ad_conversions on media.email = ad.email 
CTR = no. of conversions/impressions

# media 
# join exposures, media_customers and conversions, their own user base better 
# purchasing power of their customers across different age groups.. 

# given that the customer who bought some product is our subscriber, which type of ad_campaign generated the best revenue for the advertisers. 


# PoC1 
Advertisers Postgres Container
Media Postgres Container 
Intel SGX Edgeless DB Instance -> Result + Noise(To make private)

#PoC2 
Advertisers and Media have separate snowflake accounts on same cloud provider . 
We will generate the data loading/giving privelege/data sharing sql. 
We will give udfs for query matching/query validation/query running along with privacy measures 
Results will be visible to whoever is eligible according to contract 


