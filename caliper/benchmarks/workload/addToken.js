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
		const randID= Math.floor(Math.random()*60000)
        //const assetID = this.txIndex.toString();
        //const assetID=randID.toString()+`_${this.workerIndex}_`+this.txIndex.toString();
		const assetID=randID.toString()+`_${this.workerIndex}_${this.txIndex}`;
        ids.push(assetID)
        console.log(`Worker ${this.workerIndex}: Creating asset ${assetID}`);
        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'TokenTransaction',
            invokerIdentity: 'peer0.adv0.advnet.com',
            contractArguments: [assetID,'rand','cmp12','2','2022-02-04T00:00:01'],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
    async cleanupWorkloadModule() {
        // for (let i=0; i<ids.length; i++) {
        //     const assetID = this.ids[i];
        //     console.log(`Worker ${this.workerIndex}: Deleting asset ${assetID}`);
        //     const request = {
        //         contractId: this.roundArguments.contractId,
        //         contractFunction: 'DeleteToken',
        //         invokerIdentity: 'peer0.adv0.advnet.com',
        //         contractArguments: [assetID],
        //         readOnly: false
        //     };

        //     await this.sutAdapter.sendRequests(request);
            
        // }
        ids=[]
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