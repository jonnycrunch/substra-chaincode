package main

func registerNode(db LedgerDB, args []string) (Node, error) {
	txCreator, err := GetTxCreator(db.cc)
	if err != nil {
		return Node{}, err
	}

	node := Node{}
	node.ID = txCreator

	// Not using db.Add because we need to handle conflict as silent event without errors
	exists, err := db.KeyExists(node.ID)
	if err != nil {
		return Node{}, err
	}

	if exists {
		return node, nil
	}

	err = db.Put(node.ID, node)
	if err != nil {
		return Node{}, err
	}

	err = db.CreateIndex("node~key", []string{"node", node.ID})
	if err != nil {
		return Node{}, err
	}

	return node, nil
}

func queryNodes(db LedgerDB, args []string) (resp []Node, err error) {
	elementsKeys, err := db.GetIndexKeys("node~key", []string{"node"})
	if err != nil {
		return nil, err
	}

	nodes := []Node{}
	for _, key := range elementsKeys {
		node, err := db.GetNode(key)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}