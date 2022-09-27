package client_test

// You MUST NOT change these default imports.  ANY additional imports it will
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports. Normally, you will want to avoid underscore imports
	// unless you know exactly what you are doing. You can read more about
	// underscore imports here: https://golangdocs.com/blank-identifier-in-golang
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect(). You can read more
	// about dot imports here:
	// https://stackoverflow.com/questions/6478962/what-does-the-dot-or-period-in-a-go-import-statement-do
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	// The client implementation is intentionally defined in a different package.
	// This forces us to follow best practice and write tests that only rely on
	// client API that is exported from the client package, and avoid relying on
	// implementation details private to the client package.
	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	// We are using 2 libraries to help us write readable and maintainable tests:
	//
	// (1) Ginkgo, a Behavior Driven Development (BDD) testing framework that
	//             makes it easy to write expressive specs that describe the
	//             behavior of your code in an organized manner; and
	//
	// (2) Gomega, an assertion/matcher library that allows us to write individual
	//             assertion statements in tests that read more like natural
	//             language. For example "Expect(ACTUAL).To(Equal(EXPECTED))".
	//
	// In the Ginko framework, a test case signals failure by calling Ginkgo’s
	// Fail(description string) function. However, we are using the Gomega library
	// to execute our assertion statements. When a Gomega assertion fails, Gomega
	// calls a GomegaFailHandler, which is a function that must be provided using
	// gomega.RegisterFailHandler(). Here, we pass Ginko's Fail() function to
	// Gomega so that Gomega can report failed assertions to the Ginko test
	// framework, which can take the appropriate action when a test fails.
	//
	// This is the sole connection point between Ginkgo and Gomega.
	RegisterFailHandler(Fail)

	RunSpecs(t, "Client Tests")
}

// ================================================
// Here are some optional global variables that can be used throughout the test
// suite to make the tests more readable and maintainable than defining these
// values in each test. You can add more variables here if you want and think
// they will help keep your code clean!
// ================================================
const someFilename = "file1.txt"
const someOtherFilename = "file2.txt"
const nonExistentFilename = "thisFileDoesNotExist.txt"

const aliceUsername = "Alice"
const alicePassword = "AlicePassword"
const bobUsername = "Bob"
const bobPassword = "BobPassword"
const nilufarUsername = "Nilufar"
const nilufarPassword = "NilufarPassword"
const olgaUsername = "Olga"
const olgaPassword = "OlgaPassword"
const marcoUsername = "Marco"
const marcoPassword = "MarcoPassword"

const nonExistentUsername = "NonExistentUser"

var alice *client.User
var bob *client.User
var nilufar *client.User
var olga *client.User
var marco *client.User

var someFileContent []byte
var someShortFileContent []byte
var someLongFileContent []byte

// ================================================
// The top level Describe() contains all tests in
// this test suite in nested Describe() blocks.
// ================================================

/* ======================= Helper functions ============================ */
func copyMap(originalMap map [userlib.UUID][]byte) map [userlib.UUID][]byte {
	// Create the target map
	newMap := make(map[userlib.UUID][]byte)

	// Copy from the original map to the target map
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}

func compareSlice(a []byte, b []byte) bool {
	if len(a) != len(b) {
        return false
    }
	for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
	return true
}

var _ = Describe("Client Tests", func() {
	BeforeEach(func() {
		// This top-level BeforeEach will be run before each test.
		//
		// Resets the state of Datastore and Keystore so that tests do not
		// interfere with each other.
		userlib.DatastoreClear()
		userlib.KeystoreClear()

		userlib.SymbolicDebug = false
		userlib.SymbolicVerbose = false
	})

	BeforeEach(func() {
		// This top-level BeforeEach will be run before each test.
		//
		// Byte slices cannot be constant, so this BeforeEach resets the content of
		// each global variable to a predefined value, which allows tests to rely on
		// the expected value of these variables.
		someShortFileContent = []byte("some short file content")
		someFileContent = someShortFileContent
		someLongFileContent = []byte("some LOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOONG file content")
	})


	/* ======================= InitUser, GetUser Test ============================ */
	Describe("Creating users", func() {
		// InitUser basic check
		It("should not error when creating a new user", func() {
			_, err := client.InitUser("Alice", "password")
			Expect(err).To(BeNil(), "Failed to initialized user Alice.")
		})
        
        // Duplicated Username check
		It("should error if a username is already taken by another user", func() {
			_, err1 := client.InitUser("Alice", "password")
			Expect(err1).To(BeNil(), "Failed to initialized user Alice.")
			_, err2 := client.InitUser("Alice", "password")
            Expect(err2).ToNot(BeNil(), "Duplicated username should error.")
		})

		// Empty Username check
		It("should error if a username is empty", func(){
			_, err := client.InitUser("", "")
			Expect(err).ToNot(BeNil(), "Fail to detect empty username.")
		})

		// Username is case-sensitive
		It("should not error if two username with different letter case", func(){
			_, err := client.InitUser("Alice", "password")
			_, err = client.InitUser("alice", "password")
			Expect(err).To(BeNil(), "Fail to create case-sensitive username.")
		})

		// Duplicate password is okay
		It("should error if two users have same password", func(){
			_, err := client.InitUser("Alice", "password")
			_, err = client.InitUser("Bob", "password")
			Expect(err).To(BeNil(), "Duplicate password should be allowed.")
		})

		// Basic InitUser GetUser check
		It("should not error when get a user after create a user", func() {
			alice1, err1 := client.InitUser("Alice", "password")
			Expect(err1).To(BeNil(), "Failed to initialized user Alice.")
			alice2, err2 := client.GetUser("Alice", "password")
			Expect(err2).To(BeNil(), "Failed to get user Alice.")
			Expect(alice1).To(BeEquivalentTo(alice2),
				"The User we get is not the same as initialized User",
				alice2,
				alice1)
		})

		// GetUser not exist Username check
		It("should error if a user does not exist with that username", func() {
			_, err := client.GetUser("Alice", "password")
			Expect(err).ToNot(BeNil(), "Should error if Username never initalized.")
		})

		// GetUser with exist Username but wrong password
		It("should error if username is correct but password wrong", func() {
			client.InitUser("Alice", "password")
            _, err2 := client.GetUser("Alice", "wrong_password")
			Expect(err2).ToNot(BeNil(), "should error if password is incorrect.")
		})

		// Datastore Attack (modify user data)
		It("should error if dataStore attacker changes user data", func() {
			datastore := userlib.DatastoreGetMap()
			original_datastore := copyMap(datastore)

			client.InitUser("Alice", "password")

			for k := range datastore {
				_, ok := original_datastore[k]
				if !ok {
					userlib.DatastoreSet(k, []byte("qwertyuiop"))
				}
			}

            _, err := client.GetUser("Alice", "password")
			Expect(err).ToNot(BeNil(), "should error if user data on dataStore is changed.")
		})
		
		// Datastore Attack (remove)
		It("should error if dataStore attacker remove user data", func() {
			datastore := userlib.DatastoreGetMap()
			original_datastore := copyMap(datastore)

			client.InitUser("Alice", "password")

			for k := range datastore {
				_, ok := original_datastore[k]
				if !ok {
					userlib.DatastoreDelete(k)
				}
			}

            _, err := client.GetUser("Alice", "password")
			Expect(err).ToNot(BeNil(), "should error if user data on dataStore is removed.")
		})

		// Datastore Attack (remove all)
		It("should error if dataStore attacker remove data", func() {
			client.InitUser("Alice", "password")

			userlib.DatastoreClear()

			_, err := client.GetUser("Alice", "password")
			Expect(err).ToNot(BeNil(), "should error if all data on dataStore is removed.")
		})

	})

	/* ============================ File Tests ===============================*/
	Describe("Single user storage", func() {
		var alice *client.User

		BeforeEach(func() {
			// This BeforeEach will run before each test in this Describe block.
			alice, _ = client.InitUser("Alice", "some password")
		})

		// StoreFile
		It("should upload content without erroring", func() {
			content := []byte("This is a test")
			err := alice.StoreFile("file1", content)
			Expect(err).To(BeNil(), "Failed to upload content to a file", err)
		})

		// LoadFile nonExistentFile
		It("should error when trying to load a file that does not exist", func() {
			_, err := alice.LoadFile(nonExistentFilename)
			Expect(err).ToNot(BeNil(), "Was able to load a non-existent file without error.")
		})

		// StoreFile & LoadFile
		It("should download the expected content that was previously uploaded", func() {
			uploadedContent := []byte("This is a test")
			err := alice.StoreFile(someFilename, uploadedContent)
			Expect(err).To(BeNil(), "Failed to store a file", err)
			downloadedContent, err := alice.LoadFile(someFilename)
			Expect(err).To(BeNil(), "Failed to load a file", err)
			Expect(downloadedContent).To(BeEquivalentTo(uploadedContent),
				"Downloaded content is not the same as uploaded content",
				downloadedContent,
				uploadedContent)
		})

		// StoreFile called twice before LoadFile
		It("should download the latest content if storeFile called multi times on a same file name", func() {
			uploadedContent1 := []byte("This is a test")
			alice.StoreFile(someFilename, uploadedContent1)
			uploadedContent2 := []byte("This is also a test")
			err := alice.StoreFile(someFilename, uploadedContent2)
			Expect(err).To(BeNil(), "Failed to update a file", err)
			downloadedContent, err := alice.LoadFile(someFilename)
			Expect(err).To(BeNil(), "Failed to load a file", err)
			Expect(downloadedContent).To(BeEquivalentTo(uploadedContent2),
				"Downloaded content is not the same as latest uploaded content",
				downloadedContent,
				uploadedContent2)
		})

		// AppendToFile
		It("should append some content to a file", func() {
			content1 := []byte("First half of content")
			alice.StoreFile(someFilename, content1)
			content2 := []byte("Second half of content")
			err := alice.AppendToFile(someFilename, content2)
			Expect(err).To(BeNil(), "Fail to append content to a file")
			downloadedContent, err := alice.LoadFile(someFilename)
			Expect(err).To(BeNil(), "Failed to load a file", err)
			Expect(downloadedContent).To(BeEquivalentTo(append(content1, content2...)),
				"Downloaded content is not the same as appended file",
				downloadedContent,
				append(content1, content2...))
		})

		// AppendToFile nonExistentFile
		It("should error when try to append to non-exist filename", func() {
			err := alice.AppendToFile(someFilename, []byte("some content"))
			Expect(err).ToNot(BeNil(), "Was able to append content to non-exist filename")
		})

		// multi-session for one person
		It("should show latest change in all sessions for same user", func() {
			alice_2, _ := client.GetUser("Alice", "some password")
			content1 := []byte("First half of content")
			alice.StoreFile(someFilename, content1)
			content2 := []byte("Second half of content")
			err := alice.AppendToFile(someFilename, content2)
			downloadedContent, err := alice_2.LoadFile(someFilename)
			Expect(err).To(BeNil(), "Failed to load a file", err)
			Expect(downloadedContent).To(BeEquivalentTo(append(content1, content2...)),
				"Downloaded content is not the same as latest content",
				downloadedContent,
				append(content1, content2...)) 
		})

		It("should download the expected content that was previously uploaded", func() {
			uploadedContent := []byte("This is a test")
			err := alice.StoreFile(someFilename, uploadedContent)
			Expect(err).To(BeNil(), "Failed to store a file", err)
			downloadedContent, err := alice.LoadFile(someFilename)
			Expect(err).To(BeNil(), "Failed to load a file", err)
			Expect(downloadedContent).To(BeEquivalentTo(uploadedContent),
				"Downloaded content is not the same as uploaded content",
				downloadedContent,
				uploadedContent)
		})

		// Datastore Attack (modify before store) (0 point)
		It("should error if dataStore attacker changes before store", func() {
			datastore := userlib.DatastoreGetMap()
			for k := range datastore {
				userlib.DatastoreSet(k, []byte("qwertyuiop"))
			}

			uploadedContent := []byte("This is a test")
			err := alice.StoreFile(someFilename, uploadedContent)
			Expect(err).ToNot(BeNil(), "should error if dataStore attacker changes before store.")
		})

		// Datastore Attack (modify invitation+file data)
		It("should error if dataStore attacker changes file data", func() {
			datastore := userlib.DatastoreGetMap()
			original_datastore := copyMap(datastore)

			uploadedContent := []byte("This is a test")
			alice.StoreFile(someFilename, uploadedContent)
			
			for k := range datastore {
				_, ok := original_datastore[k]
				if !ok {
					userlib.DatastoreSet(k, []byte("qwertyuiop"))
				}
			}

            _, err := alice.LoadFile(someFilename)
			Expect(err).ToNot(BeNil(), "should error if file on dataStore is changed.")
		})
		
		// Datastore Attack (remove)
		It("should error if dataStore attacker remove data", func() {
			datastore := userlib.DatastoreGetMap()
			original_datastore := copyMap(datastore)

			uploadedContent := []byte("This is a test")
			alice.StoreFile(someFilename, uploadedContent)
			
			for k := range datastore {
				_, ok := original_datastore[k]
				if !ok {
					userlib.DatastoreDelete(k)
				}
			}

            _, err := alice.LoadFile(someFilename)
			Expect(err).ToNot(BeNil(), "should error if data on dataStore is removed.")
		})

		// Datastore Attack (modify before append) (no point)
		It("should error if dataStore attacker remove data", func() {
			datastore := userlib.DatastoreGetMap()
			original_datastore := copyMap(datastore)

			content1 := []byte("First half of content")
			alice.StoreFile(someFilename, content1)

			for k := range datastore {
				_, ok := original_datastore[k]
				if !ok {
					userlib.DatastoreSet(k, []byte("qwertyuiop"))
				}
			}

			content2 := []byte("Second half of content")
			err := alice.AppendToFile(someFilename, content2)
			Expect(err).ToNot(BeNil(), "should error if file data on dataStore is modified.")
		})

	})

	/* ============================ Sharing Tests ===============================*/
	Describe("Sharing files", func() {

		BeforeEach(func() {
			// Initialize each user to ensure the variable has the expected value for
			// the tests in this Describe() block.
			alice, _ = client.InitUser(aliceUsername, alicePassword)
			bob, _ = client.InitUser(bobUsername, bobPassword)
			nilufar, _ = client.InitUser(nilufarUsername, nilufarPassword)
			olga, _ = client.InitUser(olgaUsername, olgaPassword)
			marco, _ = client.InitUser(marcoUsername, marcoPassword)
		})

		// create invitation with non-existing file (0 point)
		It("should error when share a non-existing file", func() {
			_, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).ToNot(BeNil(), "Alice shoudn't be able to share a non-existing file.")
		})

		// accept non-existing invitation (0 point)
		It("should error when accept a non-existing invitation", func() {
			err := bob.AcceptInvitation(aliceUsername, userlib.UUIDNew(), someOtherFilename)
			Expect(err).ToNot(BeNil(), "Bob shoudn't be able to accept a non-existing invitation.")
		})

		// dataStore attack: accept modified invitation (0)
		It("should error when accept a modified invitation", func() {
			shareFileInfoPtr, _ := alice.CreateInvitation(someFilename, bobUsername)
			userlib.DatastoreSet(shareFileInfoPtr, []byte("qwertyuiop"))
			err := bob.AcceptInvitation(aliceUsername, userlib.UUIDNew(), someOtherFilename)
			Expect(err).ToNot(BeNil(), "Bob shoudn't be able to accept a modified invitation.")
		})

		// shoulf error if accept invitation with existing filename (14)
		It("should error when accept a non-existing invitation", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, _ := alice.CreateInvitation(someFilename, bobUsername)
			bob.StoreFile(someOtherFilename, someShortFileContent)
			err := bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Bob shoudn't be able to accept invitation with existing file name.")
		})

		// should error if revoke before accept invitation (24)
		It("should error when accept an invitation after revoke", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, _ := alice.CreateInvitation(someFilename, bobUsername)
			alice.RevokeAccess(someFilename, bobUsername)
			err := bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Bob shoudn't be able to accept invitation after revoke.")
		})

		// A CreateInvitation, B AcceptInvitation, B LoadFile
		It("should share a file without erroring", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")

			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")

			downloadedContent, err := bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Bob downloaded was not the same as what Alice uploaded.")
		})

		// A CreateInvitation, B AcceptInvitation, A Storefile, B LoadFile
		It("should share a file and able to update it later", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			newContent := []byte("This is also a test")
			alice.StoreFile(someFilename, newContent)

			downloadedContent, err := bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(newContent),
				"The file contents that Bob downloaded was not the same as what Alice updated.")
		})

		// A CreateInvitation, B AcceptInvitation, B Appendfile, A LoadFile
		It("should share a file and recipient should be able to update it later", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, _ := alice.CreateInvitation(someFilename, bobUsername)
			bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			newContent := []byte("This is also a test")
			bob.AppendToFile(someOtherFilename, newContent)

			downloadedContent, _ := alice.LoadFile(someFilename)
			Expect(downloadedContent).To(BeEquivalentTo(append(someShortFileContent, newContent...)),
				"The file contents that Bob downloaded was not the same as what Alice updated.")
		})

		// Steal invitation from original recipient
		It("should error if invitation send to a third person", func() {
			alice.StoreFile(someFilename, someShortFileContent)

			_, err := alice.CreateInvitation(someFilename, "non-Username")
			Expect(err).ToNot(BeNil(), "Should error if invitation is sent to non-exsiting user.")
			
			shareFileInfoPtr, _ := alice.CreateInvitation(someFilename, bobUsername)
			err = nilufar.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Should error if invitation is not intended to send to Nilufar.")
		})

		// Recognize wrong/non-existing sender
		It("should error if invitation send to a third person", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, _ := alice.CreateInvitation(someFilename, bobUsername)

			err := bob.AcceptInvitation(nilufarUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Should error if recipient recognize wrong sender.")
			
			err = bob.AcceptInvitation("non-Username", shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Should error if recipient recognize non-existing sender.")
		})

		// Test revoke:
		// A CreateInvitation, B AcceptInvitation         A
		// B CreateInvitation, N AcceptInvitation        / \
		// A CreateInvitation, O AcceptInvitation       B   O
		// O CreateInvitation, M AcceptInvitation      /     \
		// A RevokeAccess B                           N       M
		// B LoadFile should fail
		// N LoadFile should fail
		// O LoadFile should not fail
		// M LoadFile should not fail
		It("should revoke a file recursively", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")
			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")

			shareFileInfoPtr, err = bob.CreateInvitation(someOtherFilename, nilufarUsername)
			Expect(err).To(BeNil(), "Bob failed to share a file with Nilufar.")
			err = nilufar.AcceptInvitation(bobUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Nilufar could not receive the file that Bob shared.")

			shareFileInfoPtr, err = alice.CreateInvitation(someFilename, olgaUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Olga.")
			err = olga.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Olga could not receive the file that Alice shared.")

			shareFileInfoPtr, err = olga.CreateInvitation(someOtherFilename, marcoUsername)
			Expect(err).To(BeNil(), "Olga failed to share a file with Marco.")
			err = marco.AcceptInvitation(olgaUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Marco could not receive the file that Olga shared.")

			downloadedContent, err := nilufar.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Nilufar could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Nilufar downloaded was not the same as what Alice uploaded.")
			downloadedContent, err = marco.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Marco could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Marco downloaded was not the same as what Alice uploaded.")
			
			err = alice.RevokeAccess(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice can not revoke access for Bob.")

			downloadedContent, err = bob.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Bob shouldn't be able to load a file after revoke Bob.")
			downloadedContent, err = nilufar.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Nilufar shouldn't be able to load a file after revoke Bob.")
			downloadedContent, err = olga.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Olga should be able to load a file after revoke Bob.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Olga downloaded afeter revoking Bob was not the same as what Alice uploaded.")
			downloadedContent, err = marco.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Marco should be able to load a file after revoke Bob.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Marco downloaded afeter revoking Bob was not the same as what Alice uploaded.")
		})

		// revoke non-existing file (23)
		It("should error if revoke non-existing file", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, _ := alice.CreateInvitation(someFilename, bobUsername)
			bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			err := alice.RevokeAccess("non-file", bobUsername)
			Expect(err).ToNot(BeNil(), "Should error if revoke non-existing file.")
		})

		// revoke not-shared-with user (0 point)
		It("should error if revoke not-shared-with user", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			err := alice.RevokeAccess(someFilename, bobUsername)
			Expect(err).ToNot(BeNil(), "Should error if revoke not-shared-with user.")
		})

		// revoke twice (0 point)
		It("should error if revoke twice", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, _ := alice.CreateInvitation(someFilename, bobUsername)
			bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			err := alice.RevokeAccess(someFilename, bobUsername)
			err = alice.RevokeAccess(someFilename, bobUsername)
			Expect(err).ToNot(BeNil(), "Should error if revoke twice.")
		})

		// after revoke, recipient can't load/append/share (0 point)
		It("should error if revoke non-existing recipient", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, _ := alice.CreateInvitation(someFilename, bobUsername)
			bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			alice.RevokeAccess(someFilename, bobUsername)
			_, err1 := bob.LoadFile(someFilename)
			Expect(err1).ToNot(BeNil(), "Should error if load after revoke.")
			err2 := bob.AppendToFile(someFilename, []byte("rand"))
			Expect(err2).ToNot(BeNil(), "Should error if append after revoke.")
			_, err3 := bob.CreateInvitation(someFilename, nilufarUsername)
			Expect(err3).ToNot(BeNil(), "Should error if share after revoke.")
		})

		// Datastore Attack (modify invitation before accept)
		It("should error if dataStore attacker modify invitation data", func() {
			alice.StoreFile(someFilename, someShortFileContent)

			datastore := userlib.DatastoreGetMap()
			original_datastore := copyMap(datastore)

			shareFileInfoPtr, _ := alice.CreateInvitation(someFilename, bobUsername)

			for k := range datastore {
				_, ok := original_datastore[k]
				if !ok {
					userlib.DatastoreSet(k, []byte("qwertyuiop"))
				}
			}

			err := bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Should error if invitation is modified when bob accepts it.")
		})

		// Datastore Attack (delete invitation before accept)
		It("should error if dataStore attacker delete invitation data", func() {
			alice.StoreFile(someFilename, someShortFileContent)

			datastore := userlib.DatastoreGetMap()
			original_datastore := copyMap(datastore)

			shareFileInfoPtr, _ := alice.CreateInvitation(someFilename, bobUsername)

			for k := range datastore {
				_, ok := original_datastore[k]
				if !ok {
					userlib.DatastoreDelete(k)
				}
			}

			err := bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Should error if invitation is deleted when bob accepts it.")
		})

		// Datastore Attack (modify file data before send invitation) (22)
		It("should error if dataStore attacker modify invitation data", func() {
			datastore := userlib.DatastoreGetMap()
			original_datastore := copyMap(datastore)

			alice.StoreFile(someFilename, someShortFileContent)

			for k := range datastore {
				_, ok := original_datastore[k]
				if !ok {
					userlib.DatastoreSet(k, []byte("qwertyuiop"))
				}
			}

			_, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).ToNot(BeNil(), "Should error if file is modified before sending invitation.")
		})

		// Datastore Attack (modify invitation after accept) (0 point)
		It("should error if dataStore attacker modify invitation data after accept", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, _ := alice.CreateInvitation(someFilename, bobUsername)

			datastore := userlib.DatastoreGetMap()
			original_datastore := copyMap(datastore)

			bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)

			for k, val1 := range datastore {
				val2, ok := original_datastore[k]
				if !ok || !compareSlice(val1, val2) {
					userlib.DatastoreDelete(k)
				}
			}

			_, err := bob.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Should error if invitation is modified after bob accepts it.")
		})
	})
})
