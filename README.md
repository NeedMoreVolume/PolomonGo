# PolomonGo
Simple CLI to get candlestick data from Poloniex for ETH/BTC, put it in a MongoDB and print the collection entries or put it in a Cassandra DB.

List of pairs :
    btc_etc,
    btc_eth,
    btc_ltc,
    btc_rep,
    btc_xmr,
    btc_xrp,

List of chart patterns:
    SMA (20 day),
    BB (20 day)
    
Features to come:
  Ichimoku Cloud,
  More pairs


To use the CLI at all, Mongo must be installed currently. Please follow the instructions on MongoDB official site (https://docs.mongodb.com/manual/administration/install-community/) to install the most current Mongo. 

To use the CLI with Cassandra, the Cassandra DB must be installed as well. To install Cassandra, please follow the instructions on the official Cassandra site (http://cassandra.apache.org/doc/latest/getting_started/installing.html). 

Once the Cassandra is installed, create a keyspace with the following command 
`CREATE KEYSPACE cluster1 WITH REPLICATION= {'class':'SimpleStrategy', 'replication_factor':1};`. 

Using this keyspace, create your tables, use the following command and edit the market pairs as necessary 
`CREATE TABLE poloniex_btc_xrp (id uuid, timestamp double, high double, low double, open double, close double, volume double, quotevolume double, weightedaverage double, PRIMARY KEY ((id), timestamp)) WITH CLUSTERING ORDER BY (timestamp DESC);`. 

Once this configuration is complete (and the connection strings match your environment for Cassandra and Mongo), the CLI will be able to put the data in both Mongo and Cassandra.
