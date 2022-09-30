'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const file= require('./data/TPOC.json');
const verifier="verifier.adv.com:9000,verifier.pub.com:9000";
const dateS=["2022-02-04T00:00:01","2022-03-11T23:59:59","2022-05-12T00:00:01","2022-06-01T00:00:01","2022-07-01T00:00:01"];
let txN=0;
let IDs=[];
let TPOCs=[];
let count=0;
class ClaimRewardWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        //this.IDs=[]
        //this.tpocsID=[];
        
        
        for (let i=0;i<this.roundArguments.tx;i++){
            let TPOC_Set=[];

            let max=this.roundArguments.max;
            let min=this.roundArguments.min;
         
            const n=Math.floor(Math.random() * (max - min) ) + min;
            
            for(let j=0;j<n;j++){
               TPOC_Set.push(file.TPOC[j]);
            }

            const tpocClaim=TPOC_Set.join("RWRD");
            TPOCs.push(tpocClaim)
        }

        const campaignID=`${this.workerIndex}_Campaign`;
        console.log(`Worker ${this.workerIndex}: Creating asset ${campaignID}`);
        const requestC={
            contractId: this.roundArguments.contractId,
            contractFunction: 'CreateCampaign',
            invokerIdentity: 'peer0.adv0.advnet.com',
            contractArguments: [campaignID,'advtestID','pubtestID','campaign',verifier,'1','2022-01-01T00:00:01','2022-09-01T23:59:59'],
            readOnly: false
        };
        await  this.sutAdapter.sendRequests(requestC);

        for(let i=0;i<this.roundArguments.assets;i++){
            const randDate=dateS[Math.floor(Math.random()*dateS.length)];
            const connection=Math.floor(Math.random()*15) + 5
            const ctime=connection.toString();
            const request={
                contractId: this.roundArguments.contractId,
                contractFunction: 'TokenTransaction',
                invokerIdentity: 'peer0.adv0.advnet.com',
                contractArguments: [file.TPOC[i],file.TPOC_D[i],campaignID,ctime,randDate],
                readOnly: false

            };
            await this.sutAdapter.sendRequests(request);

            if(i%25==0){
                console.log(`Worker ${this.workerIndex}: Token # ${i} created`);
            }
        }
        
    }

    async submitTransaction() {
       
        const claimID=`${this.workerIndex}_${txN}`+ Math.floor(Math.random()*90000).toString();
        IDs.push(claimID);
        txN++;
        

        //console.log(tpocClaim);

        const claimRequest={
            contractId: this.roundArguments.contractId,
            contractFunction: 'ClaimReward',
            invokerIdentity: 'peer0.adv0.advnet.com',
            contractArguments: [`${this.workerIndex}_Campaign`,claimID,TPOCs[count],'timestampPlaceholder'],
            readOnly: false
        }
        await this.sutAdapter.sendRequests(claimRequest);
        count++;
    }

    async cleanupWorkloadModule() {
        const campaignID=`${this.workerIndex}_Campaign`;
        console.log(`Worker ${this.workerIndex}: Deleting asset ${campaignID}`);
        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'DeleteCampaign',
            invokerIdentity: 'peer0.adv0.advnet.com',
            contractArguments: [campaignID],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
        console.log('Deleting Tokens')
        for (let i=0; i<this.roundArguments.assets; i++) {
            const assetID = file.TPOC[i];
            console.log(`${i}_Worker ${this.workerIndex}: Deleting asset ${assetID}`);
            const requestC = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'DeleteToken',
                invokerIdentity: 'peer0.adv0.advnet.com',
                contractArguments: [assetID],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(requestC);
        }
        console.log(IDs.length)
        for (let i=0; i<IDs.length; i++) {
            const assetID = IDs[i];
            console.log(`${i}: Worker ${this.workerIndex}: Deleting asset ${assetID}`);
            const requestD = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'DeleteReward',
                invokerIdentity: 'peer0.adv0.advnet.com',
                contractArguments: [`${this.workerIndex}_Campaign`,assetID],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(requestD);
        }

        IDs=[];
        txN=0;
        count=0;
        // console.log(IDs.length)
        // console.log(txN)
    }
}

function createWorkloadModule() {
    return new ClaimRewardWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;