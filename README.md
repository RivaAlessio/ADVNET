
# ADVNET

## Running and working on the netowrk
</br>

**Script permissions**
```
sudo chmod 755 main.sh
sudo chmod 755 settings.sh
sudo chmod -R 755 scripts/
```
**Run the network**
</br> Fresh restart of the network
```
./main.sh network restart
```

<!--# </br>Available commands-->
</br>**Perform Claim Reward**
```
./main.sh chaincode testing-protocol
```
</br>**Down the network**
```
./main.sh network down
```
<!--
</br>**Generate proof from campaign verifier**
```
./main.sh chaincode proof
```
</br>**Generate proof and 2 TPoC from campaign verifier**
```
./main.sh chaincode poctpoc
```
**Stop the network**
```
./main.sh network down
```
-->

`/chaincode` directory contains all chaincode functions grouped, in `/cc` directory each function is separated to belong to its chaincode namespace

</br> :warning: :warning:
The current setup uses the same crypto parameters for each campaign, in a real implementation each different verifier should adopts different parameters for each campaign :warning: :warning:
# </br> Hyperledger Caliper Test
The current system was tested over an Intel processor equipped with 8 cores, 16 Threads. Clock speed running from 3.6 GHz to a maximum boost speed of 5.0 GHz. 16 GB of ram DDR4 3000 MHz.
* A test result example can be found in ./caliper/ folder

**Initialize caliper**
```
./main.sh caliper init
```
**Run Benchmarks** 

:warning: perform a network restart before each test to clean the ledger from previous test

Command  | Test
------------- | -------------
3tpoc  | 10000 requests at 500(tps) send rate, claim 3 Tokens at request
7tpoc  | 10000 requests at 500(tps) send rate, claim 7 Tokens at request
12tpoc  | 10000 requests at 500(tps) send rate, claim 12 Tokens at request
18tpoc  | 10000 requests at 500(tps) send rate, claim 18 Tokens at request
25tpoc  | 10000 requests at 500(tps) send rate, claim 25 Tokens at request
20000tx  | 10000 requests at 500(tps) send rate, claim 3 Tokens at request, ledger initialized with 20000 transactions
35000tx  | 10000 requests at 500(tps) send rate, claim 3 Tokens at request, ledger initialized with 35000 transactions
50000tx  | 10000 requests at 500(tps) send rate, claim 3 Tokens at request, ledger initialized with 50000 transactions
```
./main.sh caliper launch-"insert command here"
```
***example***
```
./main.sh caliper launch-3tpoc
```
**Unbind caliper and delete modules**
```
./main.sh caliper clear
```
