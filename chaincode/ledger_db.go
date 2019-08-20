package main

import (
	"chaincode/errors"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// LedgerDB to access the chaincode database during the lifetime of a SmartContract
type LedgerDB struct {
	cc shim.ChaincodeStubInterface
}

// NewLedgerDB create a new db to access the chaincode during a SmartContract
func NewLedgerDB(stub shim.ChaincodeStubInterface) LedgerDB {
	db := LedgerDB{
		cc: stub,
	}
	return db
}

// ----------------------------------------------
// Low-level functions to handle asset structs
// ----------------------------------------------

// Get retrieves an object stored in the chaincode db and set the input object value
func (db *LedgerDB) Get(key string, object interface{}) error {
	buff, err := db.cc.GetState(key)
	if err != nil {
		return errors.NotFound(err)
	}
	if buff == nil {
		return errors.NotFound()
	}
	if err = json.Unmarshal(buff, &object); err != nil {
		return err
	}
	return nil
}

// KeyExists checks if a key is stored in the chaincode db
func (db *LedgerDB) KeyExists(key string) (bool, error) {
	buff, err := db.cc.GetState(key)
	return buff != nil, err
}

// Put stores an object in the chaincode db, if the object already exists it is replaced
func (db *LedgerDB) Put(key string, object interface{}) error {
	buff, _ := json.Marshal(object)
	if err := db.cc.PutState(key, buff); err != nil {
		return err
	}
	return nil
}

// Add stores an object in the chaincode db, it fails if the object already exists
func (db *LedgerDB) Add(key string, object interface{}) error {
	ok, err := db.KeyExists(key)
	if err != nil {
		return err
	}
	if ok {
		return errors.Conflict("struct already exists (tkey: %s)", key)
	}
	return db.Put(key, object)
}

// ----------------------------------------------
// Low-level functions to handle indexes
// ----------------------------------------------

// CreateIndex adds a new composite key to the chaincode db
func (db *LedgerDB) CreateIndex(index string, attributes []string) error {
	compositeKey, err := db.cc.CreateCompositeKey(index, attributes)
	if err != nil {
		return fmt.Errorf("cannot create index %s: %s", index, err.Error())
	}
	value := []byte{0x00}
	if err = db.cc.PutState(compositeKey, value); err != nil {
		return fmt.Errorf("cannot create index %s: %s", index, err.Error())
	}
	return nil
}

// DeleteIndex deletes a composite key in the chaincode db
func (db *LedgerDB) DeleteIndex(index string, attributes []string) error {
	compositeKey, err := db.cc.CreateCompositeKey(index, attributes)
	if err != nil {
		return err
	}
	if err = db.cc.DelState(compositeKey); err != nil {
		return err
	}
	return nil
}

// UpdateIndex updates an existing composite key in the chaincode db
func (db *LedgerDB) UpdateIndex(index string, oldAttributes []string, newAttribues []string) error {
	if err := db.DeleteIndex(index, oldAttributes); err != nil {
		return err
	}
	if err := db.CreateIndex(index, newAttribues); err != nil {
		return err
	}
	return nil
}

// GetIndexKeys returns keys matching composite key values from the chaincode db
func (db *LedgerDB) GetIndexKeys(index string, attributes []string) ([]string, error) {
	keys := make([]string, 0)
	iterator, err := db.cc.GetStateByPartialCompositeKey(index, attributes)
	if err != nil {
		return nil, fmt.Errorf("get index %s failed: %s", index, err.Error())
	}
	defer iterator.Close()
	for iterator.HasNext() {
		compositeKey, err := iterator.Next()
		if err != nil {
			return nil, err
		}
		_, keyParts, err := db.cc.SplitCompositeKey(compositeKey.Key)
		if err != nil {
			return nil, fmt.Errorf("get index %s failed: cannot split key %s: %s", index, compositeKey.Key, err.Error())
		}
		keys = append(keys, keyParts[len(keyParts)-1])
	}
	return keys, nil
}

// ----------------------------------------------
// High-level functions
// ----------------------------------------------

// GetAlgo fetches an Algo from the ledger using its unique key
func (db *LedgerDB) GetAlgo(key string) (Algo, error) {
	algo := Algo{}
	if err := db.Get(key, &algo); err != nil {
		return algo, err
	}
	if algo.AssetType != AlgoType {
		return algo, errors.NotFound("algo %s not found", key)
	}
	return algo, nil
}

// GetObjective fetches an Objective from the ledger using its unique key
func (db *LedgerDB) GetObjective(key string) (Objective, error) {
	objective := Objective{}
	if err := db.Get(key, &objective); err != nil {
		return objective, err
	}
	if objective.AssetType != ObjectiveType {
		return objective, errors.NotFound("objective %s not found", key)
	}
	return objective, nil
}

// GetDataManager fetches a DataManager from the ledger using its unique key
func (db *LedgerDB) GetDataManager(key string) (DataManager, error) {
	dataManager := DataManager{}
	if err := db.Get(key, &dataManager); err != nil {
		return dataManager, err
	}
	if dataManager.AssetType != DataManagerType {
		return dataManager, errors.NotFound("dataManager %s not found", key)
	}
	return dataManager, nil
}

// GetDataSample fetches a DataSample from the ledger using its unique key
func (db *LedgerDB) GetDataSample(key string) (DataSample, error) {
	dataSample := DataSample{}
	if err := db.Get(key, &dataSample); err != nil {
		return dataSample, err
	}
	if dataSample.AssetType != DataSampleType {
		return dataSample, errors.NotFound("dataSample %s not found", key)
	}
	return dataSample, nil
}

// GetTraintuple fetches a Traintuple from the ledger using its unique key
func (db *LedgerDB) GetTraintuple(key string) (Traintuple, error) {
	traintuple := Traintuple{}
	if err := db.Get(key, &traintuple); err != nil {
		return traintuple, err
	}
	if traintuple.AssetType != TraintupleType {
		return traintuple, errors.NotFound("traintuple %s not found", key)
	}
	return traintuple, nil
}

// GetTesttuple fetches a Testtuple from the ledger using its unique key
func (db *LedgerDB) GetTesttuple(key string) (Testtuple, error) {
	testtuple := Testtuple{}
	if err := db.Get(key, &testtuple); err != nil {
		return testtuple, err
	}
	if testtuple.AssetType != TesttupleType {
		return testtuple, errors.NotFound("testtuple %s not found", key)
	}
	return testtuple, nil
}
