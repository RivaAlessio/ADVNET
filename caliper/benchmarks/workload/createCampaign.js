'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
let ids=[]

/**
 * Workload module for the benchmark round.
 */
class MyWorkloadCreate extends WorkloadModuleBase {
    /**
     * Initializes the workload module instance.
     */
    constructor() {
        super();
        this.txIndex = 0;
        
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
     async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        ids = [];
    }
    async submitTransaction() {
        this.txIndex++;
        //let carNumber = 'Client' + this.workerIndex + '_CAR' + this.txIndex.toString();
		const randID= Math.floor(Math.random()*1000)
        //const assetID = this.txIndex.toString();
        //const assetID=randID.toString()+`_${this.workerIndex}_`+this.txIndex.toString();
		const assetID=randID.toString()+`_${this.workerIndex}_${this.txIndex}`;
        ids.push(assetID)
        console.log(`Worker ${this.workerIndex}: Creating asset ${assetID}`);
        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'CreateCampaign',
            invokerIdentity: 'peer0.adv0.advnet.com',
            contractArguments: [assetID,'advtest','pubtest','campaign','verifierAdv,verifierPub','1','10:03:2022','10:05:2022'],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
    async cleanupWorkloadModule() {
        for (let i=0; i<ids.length; i++) {
            const assetID = this.ids[i];
            console.log(`Worker ${this.workerIndex}: Deleting asset ${assetID}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'DeleteCampaign',
                invokerIdentity: 'peer0.adv0.advnet.com',
                contractArguments: [assetID],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(request);
            ids=[]
        }
    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new MyWorkloadCreate();
}

module.exports.createWorkloadModule = createWorkloadModule;