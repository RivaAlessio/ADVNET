
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

# </br>Available commands
</br>**Perform Claim Reward**
```
./main.sh chaincode testing-protocol
```
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
</br> :warning: :warning:
The current setup uses the same crypto parameters for each campaign, in a real implementation each different verifier should adopts different parameters for each campaign :warning: :warning:
# </br> Hyperledger Caliper Test
**Initialize caliper**
```
./main.sh caliper init
```
**Run Benchmarks** CreateCampaign / ReadCampaign / Generate Proof / Claim Reward workloads
```
./main.sh caliper launch
```
**Unbind caliper and delete modules**
```
./main.sh caliper clear
```