package fpGrowth

type ItemSet map[int][]int

type Pattern struct {
	Items     []int
	Frequency uint
}

type FPNode struct {
	Item      int
	Frequency uint
	Link      *FPNode
	Parent    *FPNode
	Children  []*FPNode
}

func NewFPNode(item int, parent *FPNode) *FPNode {
	return &FPNode{Item: item, Parent: parent, Frequency: 1}
}

type FPTree struct {
	Root                   *FPNode
	HeaderTable            map[int]*FPNode
	MinimumSupportTreshold uint
}

func NewFPTree(transactions ItemSet, minimumSupportThreshold uint) *FPTree {
	fp := &FPTree{
		MinimumSupportTreshold: minimumSupportThreshold,
		HeaderTable:            make(map[int]*FPNode),
		Root:                   NewFPNode(0, nil),
	}

	// Find support for each item
	frequencyByItem := make(map[int]uint)
	for _, t := range transactions {
		for _, i := range t {
			frequencyByItem[i]++
		}
	}

	// Discard those not meeting the threshold
	for k, f := range frequencyByItem {
		if f < minimumSupportThreshold {
			delete(frequencyByItem, k)
		}
	}

	// Sort decending by frequency by support
	sortedFrequencies := rank(frequencyByItem)

	// Construct the FP-Tree
	for _, t := range transactions { // Look through all transactions
		currentNode := fp.Root
		for _, freq := range sortedFrequencies { // and all item frequencies
			for _, item := range t {
				if freq.Item == item { // for matching items
					found := false

					// if the node exists, increase the frequency
					for index := range currentNode.Children {
						node := currentNode.Children[index]
						if node.Item == item {
							node.Frequency++
							found = true

							// and advance to the next node
							currentNode = node
							break
						}
					}

					// Otherwise add as a new child.
					if found == false {
						newChild := NewFPNode(item, currentNode)
						currentNode.Children = append(currentNode.Children, newChild)

						// and set the new node as current
						currentNode = newChild

						// Update the linked list
						if fp.HeaderTable[item] != nil {
							prev := fp.HeaderTable[item]
							for prev.Link != nil {
								prev = prev.Link
							}
							prev.Link = newChild
						} else {
							fp.HeaderTable[item] = newChild
						}
					}
					break
				}
			}
		}
	}
	return fp
}

func (fp *FPTree) Growth() []Pattern {
	if fp.IsEmpty() {
		return nil
	}
	patterns := []Pattern{}

	if containSinglePath(fp.Root) {
		currentNode := fp.Root.Children[0]
		for currentNode != nil {
			p := Pattern{Items: []int{currentNode.Item}, Frequency: currentNode.Frequency}
			patterns = append([]Pattern{p}, patterns...)

			for _, p := range patterns {
				p.Items = append([]int{currentNode.Item}, p.Items...)
				p.Frequency = currentNode.Frequency
				if len(p.Items) > 1 {
					patterns = append([]Pattern{p}, patterns...)
				}
			}

			if len(currentNode.Children) == 1 {
				currentNode = currentNode.Children[0]
			} else {
				currentNode = nil
			}
		}
	} else {
		patternChan := make(chan []Pattern)
		count := 0
		for currentItem, node := range fp.HeaderTable {
			go fp.ProcessNode(patternChan, currentItem, node)
			count++
		}
		for count > 0 {
			patterns = append(patterns, <-patternChan...)
			count--
		}
	}
	return patterns
}

func (fp *FPTree) ProcessNode(patternChan chan []Pattern, currentItem int, node *FPNode) {
	transactionID := 1
	conditionalPatternBase := []Pattern{}
	startingNode := node
	for startingNode != nil {
		currentNode := startingNode.Parent
		if currentNode.Parent != nil {
			transformedPrefixPath := Pattern{Frequency: startingNode.Frequency}
			for currentNode.Parent != nil {
				transformedPrefixPath.Items = append(transformedPrefixPath.Items, currentNode.Item)
				currentNode = currentNode.Parent
			}
			if len(transformedPrefixPath.Items) > 1 {
				conditionalPatternBase = append(conditionalPatternBase, transformedPrefixPath)
			}
		}
		startingNode = startingNode.Link
	}

	conditionalTransactions := ItemSet{}

	for _, transformedPrefixPath := range conditionalPatternBase {
		transaction := []int{}
		transaction = append(transaction, transformedPrefixPath.Items...)

		for i := uint(0); i < transformedPrefixPath.Frequency; i++ {
			conditionalTransactions[transactionID] = transaction
			transactionID++
		}
	}

	conditionalFPTree := NewFPTree(conditionalTransactions, fp.MinimumSupportTreshold)
	conditionalPatterns := conditionalFPTree.Growth()

	currentFrequency := uint(0)

	currentPatterns := []Pattern{}
	fpNode := node
	for fpNode != nil {
		currentFrequency += fpNode.Frequency
		fpNode = fpNode.Link
	}

	for _, pattern := range conditionalPatterns {
		pattern.Items = append(pattern.Items, currentItem)
		if len(pattern.Items) > 1 {
			currentPatterns = append(currentPatterns, pattern)
		}
	}
	patternChan <- currentPatterns
}

func (fp *FPTree) IsEmpty() bool {
	if len(fp.Root.Children) > 0 {
		return false
	}
	return true
}

func containSinglePath(fpnode *FPNode) bool {
	if len(fpnode.Children) == 0 {
		return true
	} else if len(fpnode.Children) > 1 {
		return false
	}
	return containSinglePath(fpnode.Children[0])
}
