Just a toy project to work with putting data on an Ethereum blockchain. 

The idea is that you can pass in an md5 (or anything really) of a file to ethmd5, which will then save it on an Ethereum blockchain as a transaction. You can then use 
the transaction hash to grab the data and verify your file. Since the data is in the blockchain it is immutable, hopefully avoiding people jacking around with
your md5 file provided with your app for verification. Of course, you need to make the transaction hash publicly available somewhere so you still have the risk 
of someone changing it which is why this is a toy project :)

You'll need to move `sample-config.json` to `conf.json` and fill out the values. Since it is a toy you'll need your Ethereum account file in the same
directory as this app. Once you have that you can run:

`ethmd5 113e84c52c0b510893c74809dc2a83b7`
