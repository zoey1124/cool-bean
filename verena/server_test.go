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
		UUID, _ = GetUUID(aliceUsername, someFilename)
		someFileContent = []byte("some file content")
		// Generate a Merkle Tree
		var list []mt.Content
		list = append(list, LeafContent{c: someFileContent})
		someMerkleTree, _ = mt.NewTree(list)
		// Generate a FileObject
		fileObject := FileObject{Plaintext: string(someFileContent),
			MerkleTree: someMerkleTree,
			Versions:   nil}
		// Put UUID -> FileObject in DataStore
		DataStore[UUID] = fileObject
	})

	/* ============================ _loadFile Tests ================================= */
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

	/* =========================== _storeFile Tests ================================= */
	Describe("StoreFile", func() {
		It("should not error when store a file content for the first time", func() {
			bob := "Bob"
			someFileName2 := "file2.txt"
			content2 := "This is another file content"
			storeHashroot, _, err := _storeFile(bob, someFileName2, content2)
			Expect(err).To(BeNil(), "Fail to store file")

			// load file to check content
			loadHashroot, loadContent, err := _loadFile(bob, someFileName2)
			Expect(err).To(BeNil(), "Fail to load file content")
			Expect(loadContent).To(BeEquivalentTo(content2),
				"loaded content is not the same",
				loadContent,
				content2)
			Expect(loadHashroot).To(BeEquivalentTo(storeHashroot),
				"loaded hashroot is not the same",
				loadHashroot,
				storeHashroot)
		})

		It("should not err when store and load multiple times", func() {
			// store file for the first time
			bob := "Bob"
			someFileName2 := "file2.txt"
			content2 := "This is another file content"
			storeHashroot, _, err := _storeFile(bob, someFileName2, content2)
			Expect(err).To(BeNil(), "Fail to store file")

			// load file to check content
			loadHashroot, loadContent, err := _loadFile(bob, someFileName2)
			Expect(err).To(BeNil(), "Fail to load file content")
			Expect(loadContent).To(BeEquivalentTo(content2),
				"loaded content is not the same",
				loadContent,
				content2)
			Expect(loadHashroot).To(BeEquivalentTo(storeHashroot),
				"loaded hashroot is not the same",
				loadHashroot,
				storeHashroot)

			// update file
			content3 := "update on original file"
			storeHashroot, _, err = _storeFile(bob, someFileName2, content3)
			Expect(err).To(BeNil(), "Fail to load updated file content")
			// load updated file
			loadHashroot, loadContent, err = _loadFile(bob, someFileName2)
			Expect(err).To(BeNil(), "Fail to load updated file content")
			Expect(loadHashroot).To(BeEquivalentTo(storeHashroot),
				"load updated hashroot is not the same",
				loadHashroot,
				storeHashroot)
			Expect(loadContent).To(BeEquivalentTo(content3),
				"load updated content is not the same",
				loadContent,
				content3)
		})
	})
})
