'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const verifier="verifier.adv.com:9000,verifier.pub.com:9000";
class MyWorkloadRead extends WorkloadModuleBase {
    constructor() {
        super();
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

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
    }

    async submitTransaction() {
        const myArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'GenerateProof',
            invokerIdentity: 'peer0.adv0.advnet.com',
            contractArguments: [`${this.workerIndex}_Campaign`],
            readOnly: true
        };

        await this.sutAdapter.sendRequests(myArgs);
    }

    async cleanupWorkloadModule() {
        
        const assetID = `${this.workerIndex}_Campaign`;
        console.log(`Worker ${this.workerIndex}: Deleting asset ${assetID}`);
        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'DeleteCampaign',
            invokerIdentity: 'peer0.adv0.advnet.com',
            contractArguments: [assetID],
            readOnly: false
        }
        await this.sutAdapter.sendRequests(request);
        
    }
}

function createWorkloadModule() {
    return new MyWorkloadRead();
}

module.exports.createWorkloadModule = createWorkloadModule;
