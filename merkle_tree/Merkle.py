import hashlib
from typing import List


# This is a toy API of merkle tree
# It only support for building a merkle tree given a list of nodes
# Need to implement update merkle tree

class Node:
    def __init__(self, left, right, value, content, is_copied=False) -> None:
        self.left: Node = left
        self.right: Node = right
        self.value = value
        self.content = content
        self.is_copied = is_copied
    
    @staticmethod
    def hash(val) -> str:
        return hashlib.sha256(val.encode('utf-8')).hexdigest()

    def copy(self):
        return Node(self.left, self.right, self.value, self.content, True)
    

class MerkleTree:
    def __init__(self, values: List[str]):
        self.buildTree(values)
    
    # Process leaf nodes to intermediate hash nodes, call helper function
    def buildTree(self, values: List[str]):
        leaves: List[Node] = [Node(left=None, right=None, value=Node.hash(e), content=e) for e in values]
        if len(leaves) % 2 == 1:
            leaves.append(leaves[-1].copy())
        self.root: Node = self.__buildTree(leaves)
    

    def __buildTree(self, nodes: List[Node]) -> Node:
        # deal with odd nodes case
        if len(nodes) % 2 == 1:
            nodes.append(nodes[-1].copy())
        
        # base case when there are 2 nodes
        if len(nodes) == 2:
            return Node(nodes[0], nodes[1], Node.hash(nodes[0].value + nodes[1].value), nodes[0].content + "+" + nodes[1].content)

        half: int = len(nodes) // 2
        left: Node = self.__buildTree(nodes[:half])
        right: Node = self.__buildTree(nodes[half:])
        value: str = Node.hash(left.value + right.value)
        content: str = f'{left.content}+{right.content}'
        return Node(left, right, value, content)

    def getRootHash(self) -> str:
        return self.root.value

    # To be implemented
    # update rootHash(leaveNode)
    # get siblingNodes(leaveNode) -> List[Nodes]


if __name__ == "__main__":
    elems = ["A", "B", "C", "D", "E"]
    print("inputs: {}".format(elems))
    mtree = MerkleTree(elems)
    print(mtree.getRootHash())