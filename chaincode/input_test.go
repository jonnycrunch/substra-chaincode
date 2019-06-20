package main

func (dataManager *inputDataManager) createDefault() [][]byte {
	if dataManager.Name == "" {
		dataManager.Name = "liver slide"
	}
	if dataManager.OpenerHash == "" {
		dataManager.OpenerHash = dataManagerOpenerHash
	}
	if dataManager.OpenerStorageAddress == "" {
		dataManager.OpenerStorageAddress = "https://toto/dataManager/42234/opener"
	}
	if dataManager.Type == "" {
		dataManager.Type = "images"
	}
	if dataManager.DescriptionHash == "" {
		dataManager.DescriptionHash = "8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee"
	}
	if dataManager.DescriptionStorageAddress == "" {
		dataManager.DescriptionStorageAddress = "https://toto/dataManager/42234/description"
	}
	dataManager.Permissions = "all"
	args := append([][]byte{[]byte("registerDataManager")}, assetToJSON(dataManager))

	return args
}
func (dataSample *inputDataSample) createDefault() [][]byte {
	if dataSample.Hashes == "" {
		dataSample.Hashes = trainDataSampleHash1 + ", " + trainDataSampleHash2
	}
	if dataSample.DataManagerKeys == "" {
		dataSample.DataManagerKeys = dataManagerOpenerHash
	}
	if dataSample.TestOnly == "" {
		dataSample.TestOnly = "false"
	}
	args := append([][]byte{[]byte("registerDataSample")}, assetToJSON(dataSample))
	return args
}

func (objective *inputObjective) createDefault() [][]byte {
	if objective.Name == "" {
		objective.Name = "MSI classification"
	}
	if objective.DescriptionHash == "" {
		objective.DescriptionHash = objectiveDescriptionHash
	}
	if objective.DescriptionStorageAddress == "" {
		objective.DescriptionStorageAddress = "https://toto/objective/222/description"
	}
	if objective.MetricsName == "" {
		objective.MetricsName = "accuracy"
	}
	if objective.MetricsHash == "" {
		objective.MetricsHash = objectiveMetricsHash
	}
	if objective.MetricsStorageAddress == "" {
		objective.MetricsStorageAddress = objectiveMetricsStorageAddress
	}
	if objective.TestDataset == "" {
		objective.TestDataset = dataManagerOpenerHash + ":" + testDataSampleHash1 + ", " + testDataSampleHash2
	}
	objective.Permissions = "all"
	args := append([][]byte{[]byte("registerObjective")}, assetToJSON(objective))
	return args
}

func (algo *inputAlgo) createDefault() [][]byte {
	if algo.Name == "" {
		algo.Name = algoName
	}
	if algo.Hash == "" {
		algo.Hash = algoHash
	}
	if algo.StorageAddress == "" {
		algo.StorageAddress = algoStorageAddress
	}
	if algo.DescriptionHash == "" {
		algo.DescriptionHash = "e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca"
	}
	if algo.DescriptionStorageAddress == "" {
		algo.DescriptionStorageAddress = "https://toto/algo/222/description"
	}
	algo.Permissions = "all"
	args := append([][]byte{[]byte("registerAlgo")}, assetToJSON(algo))
	return args
}

func (traintuple *inputTraintuple) createDefault() [][]byte {
	if traintuple.AlgoKey == "" {
		traintuple.AlgoKey = algoHash
	}
	if traintuple.InModels == "" {
		traintuple.InModels = ""
	}
	if traintuple.ObjectiveKey == "" {
		traintuple.ObjectiveKey = objectiveDescriptionHash
	}
	if traintuple.DataManagerKey == "" {
		traintuple.DataManagerKey = dataManagerOpenerHash
	}
	if traintuple.DataSampleKeys == "" {
		traintuple.DataSampleKeys = trainDataSampleHash1 + ", " + trainDataSampleHash2
	}
	args := append([][]byte{[]byte("createTraintuple")}, assetToJSON(traintuple))
	return args
}

func (success *inputLogSuccessTrain) createDefault() [][]byte {
	if success.Key == "" {
		success.Key = traintupleKey
	}
	if success.Log == "" {
		success.Log = "no error, ah ah ah"
	}
	if success.Perf == 0 {
		success.Perf = 0.9
	}
	if success.OutModel.Hash == "" {
		success.OutModel.Hash = modelHash
	}
	if success.OutModel.StorageAddress == "" {
		success.OutModel.StorageAddress = modelAddress
	}

	args := append([][]byte{[]byte("logSuccessTrain")}, assetToJSON(success))
	return args
}
func (success *inputLogSuccessTest) createDefault() [][]byte {
	if success.Key == "" {
		success.Key = traintupleKey
	}
	if success.Log == "" {
		success.Log = "no error, ah ah ah"
	}
	if success.Perf == 0 {
		success.Perf = 0.9
	}

	args := append([][]byte{[]byte("logSuccessTest")}, assetToJSON(success))
	return args
}
func (fail *inputLogFailTrain) createDefault() [][]byte {
	if fail.Key == "" {
		fail.Key = traintupleKey
	}
	if fail.Log == "" {
		fail.Log = "man, did it failed!"
	}

	args := append([][]byte{[]byte("logFailTrain")}, assetToJSON(fail))
	return args
}
func (fail *inputLogFailTest) createDefault() [][]byte {
	if fail.Key == "" {
		fail.Key = traintupleKey
	}
	if fail.Log == "" {
		fail.Log = "man, did it failed!"
	}

	args := append([][]byte{[]byte("logFailTest")}, assetToJSON(fail))
	return args
}
func (testtuple *inputTesttuple) createDefault() [][]byte {
	if testtuple.TraintupleKey == "" {
		testtuple.TraintupleKey = traintupleKey
	}
	args, _ := inputStructToBytes(testtuple)
	args = append([][]byte{[]byte("createTesttuple")}, assetToJSON(testtuple))
	return args
}