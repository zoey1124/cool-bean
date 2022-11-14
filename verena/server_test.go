package main

import (
	"testing"

	mt "github.com/cbergoon/merkletree"
	userlib "github.com/cs161-staff/project2-userlib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Tests")
}

// ============================================================
// Some optional global variables
// ============================================================
const someFilename = "file1.txt"
const aliceUsername = "Alice"
const alicePassword = "AlicePassword"

var UUID userlib.UUID
var someFileContent []byte
var someMerkleTree *mt.MerkleTree

var _ = Describe("Server Tests", func() {
	BeforeEach(func() {
		// This top-level BeforeEach will be run before each test.
		userlib.DatastoreClear()
		userlib.KeystoreClear()

		userlib.SymbolicDebug = false
		userlib.SymbolicVerbose = false

	})

	BeforeEach(func() {
		UUID, _ = GetUUID(aliceUsername, someFilename)
		someFileContent = []byte("some file content")
		// Generate a Merkle Tree
		var list []mt.Content
		list = append(list, Content{content: someFileContent})
		someMerkleTree, _ = mt.NewTree(list)
		// Generate a FileObject
		fileObject := FileObject{Content: string(someFileContent), MerkleTree: someMerkleTree}
		// Put UUID -> FileObject in DataStore
		DataStore[UUID] = fileObject
	})

	/* ============================ loadFile Tests ================================= */
	Describe("LoadFile", func() {
		It("should not error when load file", func() {
			// use _loadFile
			hashroot, content, err := _loadFile(aliceUsername, someFilename)
			Expect(err).To(BeNil(), "Failed to load file")
			Expect(hashroot).To(BeEquivalentTo(someMerkleTree.MerkleRoot()),
				"The hashroot is not the same",
				hashroot,
				someMerkleTree.MerkleRoot())
			Expect(content).To(BeEquivalentTo(someFileContent),
				"The content is not the same",
				content,
				someFileContent)
		})
	})
})
