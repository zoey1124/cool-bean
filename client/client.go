package client

// CS 161 Project 2

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	userlib "github.com/cs161-staff/project2-userlib"

	// Likewise, useful for debugging, etc.
	"encoding/hex"

	// Useful for string mainpulation.
	"strings"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Want to import errors.
	"errors"

	// Optional. You can remove the "_" there, but please do not touch
	// anything else within the import bracket.
	_ "strconv"
)

// This serves two purposes:
// a) It shows you some useful primitives, and
// b) it suppresses warnings for items not being imported.

// This function can be safely deleted!
func someUsefulThings() {
	// Creates a random UUID
	f := userlib.UUIDNew()
	userlib.DebugMsg("UUID as string:%v", f.String())

	// Example of writing over a byte of f
	f[0] = 10
	userlib.DebugMsg("UUID as string:%v", f.String())

	// Takes a sequence of bytes and renders as hex
	h := hex.EncodeToString([]byte("fubar"))
	userlib.DebugMsg("The hex: %v", h)

	// Marshals data into a JSON representation
	// Works well with Go structures!
	d, _ := userlib.Marshal(f)
	userlib.DebugMsg("The json data: %v", string(d))
	var g userlib.UUID
	userlib.Unmarshal(d, &g)
	userlib.DebugMsg("Unmashaled data %v", g.String())

	// errors.New(...) creates an error type!
	userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("This is an error")))

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("Key is %v, %v", pk, sk)

	// Useful for string interpolation.
	_ = fmt.Sprintf("%s_%d", "file", 1)
}
/*==========================Begin Struct Definition========================*/
// Value in FileFromSender in User struct
type FromSenderInfo struct {
	SenderName string
	UUID userlib.UUID
}

type ToRecipientInfo struct {
	FilelockEncKey []byte
	FilelockMacKey []byte
	FilelockCipherUUID userlib.UUID
}

// User is the structure definition for a user record.
type User struct {
	Username string
    Password string
	SK userlib.PKEDecKey     // for receive invitation
	SigK userlib.DSSignKey   // for verify User integrity
	FilenameToSender map[string]FromSenderInfo
	FilenameToRecipient map[string]map[string]ToRecipientInfo 
}

// Filelock store information to access a file. 
type Filelock struct {
    FileEncKey []byte
	FileMacKey []byte
	FileCipherUUID userlib.UUID
}

// File represent a file structure
type File struct {
    NumOfBlock int
	// IndexToFileBlockCipherUUID map[int]userlib.UUID
}

// FileBlock represent a content subsection in a file
type FileBlock struct {
    Content []byte
}

// Cipher represents an encrypted struct
type Cipher struct {
    CipherText []byte
	Tag []byte
}

// Cipher represents an invitation
type InvitationCipher struct {
    CiphertextEnckey []byte
	CiphertextMackey []byte
	CiphertextFileLockUUID []byte
	Signature []byte
}
/*========================End of Struct Definition============================*/


/*========================== Helper functions ==============================*/

func ByteLengthNormalize(byteArr []byte, k int) ([]byte) {
	/* 
	Return a []byte with length. If input []byte len > k, trim the byte array
    If input []byte length < k, padding with 0
	*/
    if len(byteArr) >= k {
		return byteArr[:k]
	} 
	// Padding array with zero to length of k
	n := len(byteArr)
	for i := 0; i < (k - n); i++ {
		byteArr = append(byteArr, 0)
	}
	return byteArr
}

func (userdata *User) VerifyAndDecryptInvitation(invitationCipherUUID userlib.UUID, senderName string) ([]byte, []byte, userlib.UUID, error) {
    /*
	Check invitation integrity and decrypt invitation.
	Input: 
	    userdata: user who want to verify and decrypt invitation
	    invitationCipherUUID: invitationCipherUUID 
		senderName: user who send the invitation
	Return:
		filelockEncKey: filelockEncKey in Filelock struct
		filelockMacKey: filelockMacKey in Filelock struct
		filelockCipherUUID: UUID (pointer) of a filelockCipher struct
		error: error if integrityCheck failed or other errors
	*/
	// 1 Get InvitationCipher from DataStore
	value, exist := userlib.DatastoreGet(invitationCipherUUID)
	if !exist { return nil, nil, userlib.UUIDNew(), errors.New(strings.ToTitle("Invitation Cipher does not exist.")) }
	var invitationCipher InvitationCipher
	userlib.Unmarshal(value, &invitationCipher)

	// 2 Verify InvitationCipher integrity with signature
	vk, found := userlib.KeystoreGet(senderName + "VK")
	if !found { return nil, nil, userlib.UUIDNew(), errors.New(strings.ToTitle("Sender invalid.")) }
	cipherText := append(invitationCipher.CiphertextEnckey, invitationCipher.CiphertextMackey...)
	cipherText = append(cipherText, invitationCipher.CiphertextFileLockUUID...)
	err := userlib.DSVerify(vk, cipherText, invitationCipher.Signature)
	if err!=nil { return nil, nil, userlib.UUIDNew(), err }

	// 3 Decrypt ciphertext to filelockEncKey, filelockMacKey and FilelockCipherUUID
	filelockEncKey, err1 := userlib.PKEDec((*userdata).SK, invitationCipher.CiphertextEnckey)
	filelockMacKey, err2 := userlib.PKEDec((*userdata).SK, invitationCipher.CiphertextMackey)
	plaintext, err3 := userlib.PKEDec((*userdata).SK, invitationCipher.CiphertextFileLockUUID)
	if err1!=nil { return nil, nil, userlib.UUIDNew(), err1 }
	if err2!=nil { return nil, nil, userlib.UUIDNew(), err2 }
	if err3!=nil { return nil, nil, userlib.UUIDNew(), err3 }
	var filelockCipherUUID userlib.UUID
	userlib.Unmarshal(plaintext, &filelockCipherUUID)
	return filelockEncKey, filelockMacKey, filelockCipherUUID, nil
}

func VerifyAndDecryptFilelock(filelockCipherUUID userlib.UUID, filelockEncKey []byte, filelockMacKey []byte) (*Filelock, error) {
	/* 
	Input:
		filelockCipherUUID: userlib.UUID
		filelockEncKey: []byte
		filelockMacKey: []byte
	Return: 
		&filelock: *Filelock
		error: error if integrityCheck failed or other errors
	*/
	// 1 Get FileLockCipher from DataStore
	value, exist := userlib.DatastoreGet(filelockCipherUUID)
	if !exist { return nil, errors.New(strings.ToTitle("FileLock Cipher does not exist.")) }
	var filelockCipher Cipher
	userlib.Unmarshal(value, &filelockCipher)

	// 2 Verify FileLockCipher integrity with tag. If tag changed, return error
	tag, _ := userlib.HMACEval(filelockMacKey, filelockCipher.CipherText)
	ok := userlib.HMACEqual(tag, filelockCipher.Tag)
	if !ok { return nil, errors.New(strings.ToTitle("Datastore Attacked!")) }

	// 3 Decrypt ciphertext to Filelock struct
	plaintext := userlib.SymDec(filelockEncKey, filelockCipher.CipherText)
	var filelock Filelock
	userlib.Unmarshal(plaintext, &filelock)
	return &filelock, nil
}

func VerifyAndDecryptFile(filelock *Filelock)(*File, error) {
	/* 
	Input:
		filelock: *Filelock, or a pointer to a Filelock struct
	Return:
		&file: *File, or a pointer to a File struct
		error: error if integrityCheck failed or other errors
	*/
	// 1 Get fileEncKey & fileMaccKey from filelock
	fileEncKey := (*filelock).FileEncKey
	fileMacKey := (*filelock).FileMacKey
	
	// 2 Get File from DataStore
	value, exist := userlib.DatastoreGet((*filelock).FileCipherUUID)
	if !exist { return nil, errors.New(strings.ToTitle("File Cipher does not exist.")) }
	var fileCipher Cipher
	userlib.Unmarshal(value, &fileCipher)

	// 3 Verify FileCipher integrity with tag. If tag changed, return error
	tag, _ := userlib.HMACEval(fileMacKey, fileCipher.CipherText)
	ok := userlib.HMACEqual(tag, fileCipher.Tag)
	if !ok { return nil, errors.New(strings.ToTitle("Datastore Attacked!")) }

	// 4 Decrypt ciphertext to File struct
	plaintext := userlib.SymDec(fileEncKey, fileCipher.CipherText)
	var file File
	userlib.Unmarshal(plaintext, &file)
	return &file, nil
}

func VerifyAndDecryptFileBlock(fileBlockCipherUUID userlib.UUID, fileEncKey []byte, fileMacKey []byte)(*FileBlock, error) {
	/* 
	Input:
		fileBlockCipherUUID: userlib.UUID
		fileEncKey: []byte
		fileMacKey: []byte
	Return:
		&fileBlock: *FileBlock, or a pointer to a FileBlock struct
		error: error if integrityCheck failed or other errors
	*/
	// 1 Get FileBlock from DataStore
	value, exist := userlib.DatastoreGet(fileBlockCipherUUID)
	if !exist { return nil, errors.New(strings.ToTitle("FileBlock Cipher does not exist.")) }
	var fileBlockCipher Cipher
	userlib.Unmarshal(value, &fileBlockCipher)

	// 2 Verify FileBlockCipher integrity with tag. If tag changed, return error
	tag, _ := userlib.HMACEval(fileMacKey, fileBlockCipher.CipherText)
	ok := userlib.HMACEqual(tag, fileBlockCipher.Tag)
	if !ok { return nil, errors.New(strings.ToTitle("Datastore Attacked 3!")) }

	// 3 Decrypt ciphertext to FileBlock struct
	plaintext := userlib.SymDec(fileEncKey, fileBlockCipher.CipherText)
	var fileBlock FileBlock
	userlib.Unmarshal(plaintext, &fileBlock)
	return &fileBlock, nil
}

func IntToByte(i int) ([]byte) {
	/* 
	Transform int64 i to a byte array slice. Return the byte array.
	*/
	b := []byte{byte(0xff & i), 
		        byte(0xff & (i >> 8)), 
				byte(0xff & (i >> 16)), 
				byte(0xff & (i >> 24)),
				byte(0xff & (i >> 32)),
				byte(0xff & (i >> 40)),
				byte(0xff & (i >> 48)),
				byte(0xff & (i >> 56))}
	return b
}

func CreateAndStoreFileBlock(content []byte, fileEncKey []byte, fileMacKey []byte, index int)(userlib.UUID) {
	/* 
	Create a file block with content and encrypt it into a cipher. 
	Store the cipher to dataStore. 
	Input:
		content: []byte
		fileEncKey: []byte
		fileMacKey: []byte
		index: index = file.NumOfBlock - 1
	Return:
		fileBlockCipherUUID: userlib.UUID, generate from Hash(fileEncKey||index)
	*/
	fileBlock := FileBlock{Content: content}
	iv := userlib.RandomBytes(16)
	plaintext, _ := userlib.Marshal(&fileBlock)
	ciphertext := userlib.SymEnc(fileEncKey, iv, plaintext)
	tag, _ := userlib.HMACEval(fileMacKey, ciphertext)
	fileBlockCipher := Cipher{ciphertext, tag}
	index_byte := IntToByte(index)
	fileBlockCipherUUID, _ := userlib.UUIDFromBytes(userlib.Hash(append(fileEncKey, index_byte...)))
	fileBlockCipherValue, _ := userlib.Marshal(&fileBlockCipher)
	userlib.DatastoreSet(fileBlockCipherUUID, fileBlockCipherValue)
	return fileBlockCipherUUID
}

func (userdata *User) CreateAndStoreFromSenderInfo(invitationUUID userlib.UUID, filename string, sendername string) (userlib.UUID) {
	/* 
	Generate FromSenderInfo
	Encrypt FromSenderInfo with userEncKey and userMacKey.
	Store in Datastore with UUID(H(username||filename)).
	Input:
		invitationUUID
		filename
		sendername
	Return:
	    UUID(H(username||filename))
	*/

	// Get UserEncKey and UserMacKey
	username_byte := ByteLengthNormalize([]byte(userdata.Username), 16)
	password_byte := []byte(userdata.Password)
	filename_byte := ByteLengthNormalize([]byte(filename), 16)
	generalKey := userlib.Argon2Key(password_byte, username_byte, 16)
	userEncKey, _ := userlib.HashKDF(generalKey, []byte("user encryption"))
	userEncKey = userEncKey[:16]
	userMacKey, _ := userlib.HashKDF(generalKey, []byte("user mac"))
	userMacKey = userMacKey[:16]

	// Encrypt invitationUUID
	fromSenderInfo := FromSenderInfo{SenderName: sendername, UUID: invitationUUID}
	plaintext, _ := userlib.Marshal(fromSenderInfo)
	iv := userlib.RandomBytes(16)
	ciphertext := userlib.SymEnc(userEncKey, iv, plaintext)
	tag, _ := userlib.HMACEval(userMacKey, ciphertext)
	cipher := Cipher{ciphertext, tag}

	// Generate UUID(H(username||filename))
	UUID, _ := userlib.UUIDFromBytes(userlib.Hash(append(username_byte, filename_byte...)))
	value, _ := userlib.Marshal(&cipher)
	userlib.DatastoreDelete(UUID)
	userlib.DatastoreSet(UUID, value)
	return UUID
}

func (userdata *User)VerifyAndDecryptFromSenderInfo(filename string) (*FromSenderInfo, error) {
	/*
	Generate UUID(H(username||filename)) to get FromSenderInfo object.
	Verify with UserMacKey; Decrypt with UserEncKey
	Input:
		userdata
		filename
	Return: 
		invitationUUID
		error
	*/
	// 1 Generate UUID(H(username||filename))
	username_byte := ByteLengthNormalize([]byte(userdata.Username), 16)
	password_byte := []byte(userdata.Password)
	filename_byte := ByteLengthNormalize([]byte(filename), 16)
	generalKey := userlib.Argon2Key(password_byte, username_byte, 16)
	userEncKey, _ := userlib.HashKDF(generalKey, []byte("user encryption"))
	userEncKey = userEncKey[:16]
	userMacKey, _ := userlib.HashKDF(generalKey, []byte("user mac"))
	userMacKey = userMacKey[:16]
	UUID, _ := userlib.UUIDFromBytes(userlib.Hash(append(username_byte, filename_byte...)))

	// 2 Get Cipher object from Datastore
	value, exist := userlib.DatastoreGet(UUID)
	if !exist { return nil, errors.New(strings.ToTitle("Invitation UUID Cipher does not exist.")) }
	var cipher Cipher
	userlib.Unmarshal(value, &cipher)

	// 3 Verify Integrity with tag
	tag, _ := userlib.HMACEval(userMacKey, cipher.CipherText)
	ok := userlib.HMACEqual(tag, cipher.Tag)
	if !ok { return nil, errors.New(strings.ToTitle("Datastore Attacked 4!")) }

	// 4 Decrypt to get invitationUUID
	plaintext := userlib.SymDec(userEncKey, cipher.CipherText)
	var fromSenderInfo FromSenderInfo
	userlib.Unmarshal(plaintext, &fromSenderInfo)
	return &fromSenderInfo, nil
}

func StoreFile(fileCipherUUID userlib.UUID, file *File, fileEncKey []byte, fileMacKey []byte)() {
	/* 
	Encrypt and store the file the cipher to dataStore. 
	Input:
		fileCipherUUID: userlib.UUID
		file: *File
		fileEncKey: []byte
		fileMacKey: []byte
	Return:
	*/
	plaintext, _ := userlib.Marshal(file)
	iv := userlib.RandomBytes(16)
	ciphertext := userlib.SymEnc(fileEncKey, iv, plaintext)
	tag, _ := userlib.HMACEval(fileMacKey, ciphertext)
	fileCipher := Cipher{CipherText: ciphertext, Tag: tag}
	fileCipherValue, _ := userlib.Marshal(&fileCipher)
	userlib.DatastoreDelete(fileCipherUUID)
	userlib.DatastoreSet(fileCipherUUID, fileCipherValue)
	return
}

func StoreFilelock(filelockCipherUUID userlib.UUID, filelock *Filelock, filelockEncKey []byte, filelockMacKey []byte)() {
	/* 
	Encrypt and store the filelock the cipher to dataStore. 
	Input:
		filelockCipherUUID: userlib.UUID
		filelock: *File
		filelockEncKey: []byte
		filelockMacKey: []byte
	Return:
	*/
	plaintext, _ := userlib.Marshal(filelock)
		iv := userlib.RandomBytes(16)
		ciphertext := userlib.SymEnc(filelockEncKey, iv, plaintext)
		tag, _ := userlib.HMACEval(filelockMacKey, ciphertext)
		filelockCipher := Cipher{CipherText: ciphertext, Tag: tag}
		filelockCipherValue, _ := userlib.Marshal(&filelockCipher)
		userlib.DatastoreDelete(filelockCipherUUID)
		userlib.DatastoreSet(filelockCipherUUID, filelockCipherValue)
	return
}

func (userdata *User) CreateAndStoreInvitation(invitationCipherUUID userlib.UUID, receiverName string, filelockEncKey []byte, filelockMacKey []byte, filelockCipherUUID userlib.UUID)(error) {
	/* 
	Create and store the file cipher to dataStore. 
	Input:
		userdata: user who want to verify and decrypt invitation
		invitationCipherUUID: userlib.UUID
		receiverName: string
		filelockEncKey: []byte
		filelockMacKey: []byte
		filelockCipherUUID: userlib.UUID
	Return:
		error	
	*/
	pk, found := userlib.KeystoreGet(receiverName + "PK")
	if !found { return errors.New(strings.ToTitle("Recipient invalid.")) }
	ciphertextEnckey, _ := userlib.PKEEnc(pk, filelockEncKey)
	ciphertextMackey, _ := userlib.PKEEnc(pk, filelockMacKey)
	marshaledFilelockCipherUUID, _ := userlib.Marshal(filelockCipherUUID)
	ciphertextFileLockUUID, _ := userlib.PKEEnc(pk, marshaledFilelockCipherUUID)
	ciphertext := append(ciphertextEnckey, ciphertextMackey...)
	ciphertext = append(ciphertext, ciphertextFileLockUUID...)
	signature, _ := userlib.DSSign((*userdata).SigK, ciphertext)
	invitationCipher := InvitationCipher{CiphertextEnckey: ciphertextEnckey, CiphertextMackey: ciphertextMackey, CiphertextFileLockUUID: ciphertextFileLockUUID, Signature: signature}
	invitationCipherValue, _ := userlib.Marshal(&invitationCipher)
	userlib.DatastoreDelete(invitationCipherUUID)
	userlib.DatastoreSet(invitationCipherUUID, invitationCipherValue)
	return nil
}

func (userdata *User) StoreUser()() {
	/* 
	Encrypt and store the user cipher to dataStore. 
	Input:
	Return:
	*/
	username_byte := ByteLengthNormalize([]byte((*userdata).Username), 16)
	password_byte := []byte((*userdata).Password)
	key, _ := userlib.UUIDFromBytes(username_byte)
	value, _ := userlib.DatastoreGet(key)
	var userCipher Cipher
	userlib.Unmarshal(value, &userCipher)
	plaintext, _ := userlib.Marshal(userdata)
	iv := userlib.RandomBytes(16)
	generalKey := userlib.Argon2Key(password_byte, username_byte, 16)
	userEncKey, _ := userlib.HashKDF(generalKey, []byte("user encryption"))
	userEncKey = userEncKey[:16]
	userMacKey, _ := userlib.HashKDF(generalKey, []byte("user mac"))
	userMacKey = userMacKey[:16]
	ciphertext := userlib.SymEnc(userEncKey, iv, plaintext)
	userCipher.CipherText = ciphertext
	tag, _ := userlib.HMACEval(userMacKey, ciphertext)
	userCipher.Tag = tag
	userCipherValue, _ := userlib.Marshal(&userCipher)
	userlib.DatastoreDelete(key)
	userlib.DatastoreSet(key, userCipherValue)
	return
}

func DeleteFileBlock(fileEncKey []byte, numOfBlock int) {
	/* 
	Delete all content from the file with fileEncKey
	*/
	for i := 0; i < numOfBlock; i++ {
		index_byte := IntToByte(i)
		fileBlockCipherUUID, _ := userlib.UUIDFromBytes(userlib.Hash(append(fileEncKey, index_byte...)))
		userlib.DatastoreDelete(fileBlockCipherUUID)
	}
}

/* ========================================================================== */


/* ============================ API functions============================== */
func InitUser(username string, password string) (userdataptr *User, err error) {
	// TODO: maximum length constraints on username?
	// TODO: change User UUID from UUID(username) to UUID(H(username))
	// Check Username length bigger then 0
	if len(username) == 0 {
		return nil, errors.New(strings.ToTitle("Username cannot be empty."))
	}
	// Return an error if Username is duplicated
	username_byte := ByteLengthNormalize([]byte(username), 16)
	password_byte := []byte(password)
	key, _ := userlib.UUIDFromBytes(username_byte)
	_, duplicated := userlib.DatastoreGet(key)
	if duplicated {
		return nil, errors.New(strings.ToTitle("Username is duplicated."))
	}
    
	// Create a User struct
	userdata := User{Username: username, Password: password, 
				FilenameToSender: make(map[string]FromSenderInfo),
				FilenameToRecipient: make(map[string]map[string]ToRecipientInfo)}
	// Generate key pair for public-key encryption for received invitation
	pk, sk, _ := userlib.PKEKeyGen()
	userdata.SK = sk 
	userlib.KeystoreSet(username + "PK", pk)
    // Generate key pair for Digital Digniture for received invitation
	var sigk userlib.DSSignKey
	var vk userlib.DSVerifyKey
	sigk, vk, _ = userlib.DSKeyGen()
	userdata.SigK = sigk 
	userlib.KeystoreSet(username + "VK", vk)

	// Create a Cipher struct to encrypt userdata 
    var userCipher Cipher
    var userEncKey []byte
	plaintext, _ := userlib.Marshal(&userdata)
	generalKey := userlib.Argon2Key(password_byte, username_byte, 16)
	userEncKey, _ = userlib.HashKDF(generalKey, []byte("user encryption"))
	userEncKey = userEncKey[:16]
	iv := userlib.RandomBytes(16)
	ciphertext := userlib.SymEnc(userEncKey, iv, plaintext)
	userCipher.CipherText = ciphertext 
    // Create a tag on ciphertext by HMAC 
	var userMacKey, tag []byte
	userMacKey, _ = userlib.HashKDF(generalKey, []byte("user mac"))
	userMacKey = userMacKey[:16]
	tag, _ = userlib.HMACEval(userMacKey, ciphertext)
	userCipher.Tag = tag 

	// Store User Ciper into DataStore with UUID key from username
	value, _ := userlib.Marshal(&userCipher)
    userlib.DatastoreSet(key, value)
	
	return &userdata, nil
}


/*=============================================================================*/
func GetUser(username string, password string) (userdataptr *User, err error) {
	// Get User cipher from DataStore with UUID from Username
	username_byte := ByteLengthNormalize([]byte(username), 16)
	password_byte := []byte(password)
	var key userlib.UUID
	key, _ = userlib.UUIDFromBytes(username_byte)
	value, exist := userlib.DatastoreGet(key)
	if !exist { return nil, errors.New(strings.ToTitle("Username does not exist.")) }
	var userCipher Cipher
	userlib.Unmarshal(value, &userCipher)

	// Verify userCipher integrity with tag. If tag changed, return error
	generalKey := userlib.Argon2Key(password_byte, username_byte, 16)
	var userMacKey, tag []byte
	userMacKey, _ = userlib.HashKDF(generalKey, []byte("user mac"))
	userMacKey = userMacKey[:16]
	tag, _ = userlib.HMACEval(userMacKey, userCipher.CipherText)
	ok := userlib.HMACEqual(tag, userCipher.Tag)
	if !ok { return nil, errors.New(strings.ToTitle("Datastore Attacked!")) }

	// Decrypt ciphertext to User struct
	userEncKey, _ := userlib.HashKDF(generalKey, []byte("user encryption"))
	userEncKey = userEncKey[:16]
	var userdata User
	plaintext := userlib.SymDec(userEncKey, userCipher.CipherText)
	userlib.Unmarshal(plaintext, &userdata)

	return &userdata, nil
}


/*================================= StoreFile =======================================*/
func (userdata *User) StoreFile(filename string, content []byte) (err error) {
    // 1. Reload User struct from DataStore to ensure changes made by other sessions are seen
	var user_error error
	userdata, user_error = GetUser((*userdata).Username, (*userdata).Password)
	if user_error != nil {return user_error}

	// 2. Check if filename is in filemap in User
    infoFromSender, ok := (*userdata).FilenameToSender[filename]
 	// 2.1 filename Does not exist, create file
	 if !ok {
		// 2.1.1 Create random fileEncKey, fileMacKey, filelockEncKey and filelockMacKey
		fileEncKey := userlib.RandomBytes(16)
		fileMacKey := userlib.RandomBytes(16)
		filelockEncKey := userlib.RandomBytes(16)
		filelockMacKey := userlib.RandomBytes(16)

		// 2.1.2 Create, encrpyt and store a file block to dataStore
		CreateAndStoreFileBlock(content, fileEncKey, fileMacKey, 0)

        // 2.1.3 Create, encrpyt and store a file to dataStore
		file := File{NumOfBlock: 1}
		fileCipherUUID := userlib.UUIDNew()
		StoreFile(fileCipherUUID, &file, fileEncKey, fileMacKey)

		// 2.1.4 Create, encrypt and store a filelock to datastore 
		filelock := Filelock{FileEncKey: fileEncKey, FileMacKey: fileMacKey, 
			FileCipherUUID: fileCipherUUID}
		filelockCipherUUID := userlib.UUIDNew()
		StoreFilelock(filelockCipherUUID, &filelock, filelockEncKey, filelockMacKey)

		// 2.1.5 Create and store a invitation cipher to datastore
		invitationCipherUUID := userlib.UUIDNew()
		err := userdata.CreateAndStoreInvitation(invitationCipherUUID, (*userdata).Username, filelockEncKey, filelockMacKey, filelockCipherUUID)
		if err!=nil {return err}

		// 2.1.6 Store Invitation info in User's map
		(*userdata).FilenameToSender[filename] = FromSenderInfo{SenderName: (*userdata).Username, 
											UUID: invitationCipherUUID}
		(*userdata).FilenameToRecipient[filename] = make(map[string]ToRecipientInfo)
		toRecipientInfo := ToRecipientInfo{FilelockEncKey: filelockEncKey, FilelockMacKey: filelockMacKey, 
			FilelockCipherUUID: filelockCipherUUID}
		(*userdata).FilenameToRecipient[filename][(*userdata).Username] = toRecipientInfo

		// 2.1.7 Encrypt updated User to User Cipher. Store back to Datastore
		userdata.StoreUser()

		// 2.1.8 Generate FromSenderInfo Cipher, store with UUID(H(username||filename))
		userdata.CreateAndStoreFromSenderInfo(invitationCipherUUID, filename, userdata.Username)

	// 2.2 If filename exist, update file
	} else {
		//2.2.1 Check invitation integrity and decrypt invitation to get filelock info.
		invitationCipherUUID := infoFromSender.UUID
		filelockEncKey, filelockMacKey, filelockCipherUUID, err := userdata.VerifyAndDecryptInvitation(invitationCipherUUID, infoFromSender.SenderName)
		if err!=nil {return err}

		// 2.2.2 Check filelock integrity and decrypt filelock to get fileEncKey & fileMacKey.
		filelock, err := VerifyAndDecryptFilelock(filelockCipherUUID, filelockEncKey, filelockMacKey)
		if err!=nil {return err}
		fileEncKey := (*filelock).FileEncKey
		fileMacKey := (*filelock).FileMacKey

		// 2.2.3 Check file integrity and decrypt file.
		file, err := VerifyAndDecryptFile(filelock)
		if err!=nil {return err}

		// Delete original file block
		DeleteFileBlock(fileEncKey, file.NumOfBlock)

		// 2.2.4 Create, encrpyt and store a file block to dataStore
		CreateAndStoreFileBlock(content, fileEncKey, fileMacKey, 0)

		//2.2.5 Update File with only the new file block
		(*file).NumOfBlock = 1

		//2.2.6 Encrypt and store the file back to Datastore
		plaintext, _ := userlib.Marshal(file)
		iv := userlib.RandomBytes(16)
		ciphertext := userlib.SymEnc(fileEncKey, iv, plaintext)
		tag, _ := userlib.HMACEval(fileMacKey, ciphertext)
		fileCipher := Cipher{CipherText: ciphertext, Tag: tag}
		fileCipherValue, _ := userlib.Marshal(&fileCipher)
		userlib.DatastoreDelete((*filelock).FileCipherUUID)
		userlib.DatastoreSet((*filelock).FileCipherUUID, fileCipherValue)
	}

	return 
}


/*=============================== Load File =========================================*/
func (userdata *User) LoadFile(filename string) (content []byte, err error) {
    // 1 Reload userdata to update change from other sessions
	var user_error error
	userdata, user_error = GetUser((*userdata).Username, (*userdata).Password)
	if user_error != nil {return nil, user_error}

	// 2 Check if filename is in filemap in User, return error if not
	infoFromSender, ok := (*userdata).FilenameToSender[filename]
	if !ok {
		return nil, errors.New(strings.ToTitle("File not found!"))
	}

	// 3 Check invitation integrity and decrypt invitation to get filelock info.
	invitationCipherUUID := infoFromSender.UUID
	filelockEncKey, filelockMacKey, filelockCipherUUID, err := userdata.VerifyAndDecryptInvitation(invitationCipherUUID, infoFromSender.SenderName)
	if err!=nil {return nil, err}

	// 4 Check filelock integrity and decrypt filelock to get fileEncKey & fileMacKey.
	filelock, err := VerifyAndDecryptFilelock(filelockCipherUUID, filelockEncKey, filelockMacKey)
	if err!=nil {return nil, err}
	fileEncKey := (*filelock).FileEncKey
	fileMacKey := (*filelock).FileMacKey

	// 5 Check file integrity and decrypt file.
	file, err := VerifyAndDecryptFile(filelock)
	if err!=nil {return nil, err}

	// 6 Read and concatenate all file blocks stored in file in order of index.
	for i:=0; i<(*file).NumOfBlock; i++ {
		// 6.1 Check file block integrity and decrypt file block.
		index_byte := IntToByte(i)
		fileBlockCipherUUID, _ := userlib.UUIDFromBytes(userlib.Hash(append(fileEncKey, index_byte...)))
		fileBlock, err := VerifyAndDecryptFileBlock(fileBlockCipherUUID, fileEncKey, fileMacKey)
		if err!=nil {return nil, err}

		// 6.2 Append FileBlock content to output
		content = append(content, (*fileBlock).Content...)
	}
	return content, nil
}


/*=============================== Append File ===================================*/
func (userdata *User) AppendToFile(filename string, content []byte) error {
	// 1 Check if UUID(H(username||filename)) is in Datastore, error if not
	fromSenderInfo, err := userdata.VerifyAndDecryptFromSenderInfo(filename)
	if err != nil {return err}

	// 2 Check invitation integrity and decrypt invitation to get filelock info.
	senderName := fromSenderInfo.SenderName
	invitationCipherUUID := fromSenderInfo.UUID
	filelockEncKey, filelockMacKey, filelockCipherUUID, err := userdata.VerifyAndDecryptInvitation(invitationCipherUUID, senderName)
	if err!=nil {return err}

	// 4 Check filelock integrity and decrypt filelock to get fileEncKey & fileMacKey.
	filelock, err := VerifyAndDecryptFilelock(filelockCipherUUID, filelockEncKey, filelockMacKey)
	if err!=nil {return err}
	fileEncKey := (*filelock).FileEncKey
	fileMacKey := (*filelock).FileMacKey

	// 5 Check file integrity and decrypt file.
	file, err := VerifyAndDecryptFile(filelock)
	if err!=nil {return err}

	// 6 Create, encrpyt and store a file block to dataStore
	numOfBlock := file.NumOfBlock
    CreateAndStoreFileBlock(content, fileEncKey, fileMacKey, numOfBlock)
    
	// 7 Update NumOfBlock in File; Encrypt and store the file back to Datastore
	(*file).NumOfBlock += 1
	fileCipherUUID := (*filelock).FileCipherUUID
	StoreFile(fileCipherUUID, file, fileEncKey, fileMacKey)

	return nil
}

/* ======================== Create Invitation =================================== */
func (userdata *User) CreateInvitation(filename string, recipientName string) (
	invitationPtr userlib.UUID, err error) {
	// 1. Reload User to update from other sessions
	var user_error error
	userdata, user_error = GetUser((*userdata).Username, (*userdata).Password)
	if user_error != nil {return userlib.UUIDNew(), user_error}

	// 2. Check if filename is in user's filemap
	infoFromSender, ok := userdata.FilenameToSender[filename]
	if !ok {
		return userlib.UUIDNew(), errors.New(strings.ToTitle("File not found."))
	}

	// 3. Check invitation integrity and decrypt invitation to get filelock info.
	invitationCipherUUID := infoFromSender.UUID
	filelockEncKey, filelockMacKey, filelockCipherUUID, err := userdata.VerifyAndDecryptInvitation(invitationCipherUUID, infoFromSender.SenderName)
	if err!=nil {return userlib.UUIDNew(), err}
	
	// 4. Check if user is the owner of the file
	var toRecipientInfo ToRecipientInfo
	newInvitationCipherUUID := userlib.UUIDNew()
	if infoFromSender.SenderName == (*userdata).Username {
		// 4.1. If user is the owner, create new filelock and invitation cipher.
		// 4.1.1. Check received filelock integrity and decrypt filelock to get fileEncKey & fileMacKey.
		filelock, err := VerifyAndDecryptFilelock(filelockCipherUUID, filelockEncKey, filelockMacKey)
		if err!=nil {return userlib.UUIDNew(), err}
		fileEncKey := (*filelock).FileEncKey
		fileMacKey := (*filelock).FileMacKey
		fileCipherUUID := (*filelock).FileCipherUUID

		// 4.1.2. Create new FilelockEncKey, FilelockMacKey
		newFilelockEncKey := userlib.RandomBytes(16)
		newFilelockMacKey := userlib.RandomBytes(16)

		// 4.1.3. Create, encrypt and store a filelock to datastore 
		*filelock = Filelock{FileEncKey: fileEncKey, FileMacKey: fileMacKey, 
			FileCipherUUID: fileCipherUUID}
		newFilelockCipherUUID := userlib.UUIDNew()
		StoreFilelock(newFilelockCipherUUID, filelock, newFilelockEncKey, newFilelockMacKey)

		// 4.1.4. Create and store a invitation cipher to datastore
		err = userdata.CreateAndStoreInvitation(newInvitationCipherUUID, recipientName, newFilelockEncKey, newFilelockMacKey, newFilelockCipherUUID)
		if err!=nil {return userlib.UUIDNew(), err}

		toRecipientInfo = ToRecipientInfo{FilelockEncKey: newFilelockEncKey, FilelockMacKey: newFilelockMacKey, 
			FilelockCipherUUID: newFilelockCipherUUID}
			
	}  else {
		// 4.2 If user not the owner, create invitation cipher containing existing filelock info.
		// 4.2.1. Create and store a invitation cipher to datastore
		err := userdata.CreateAndStoreInvitation(newInvitationCipherUUID, recipientName, filelockEncKey, filelockMacKey, filelockCipherUUID)
		if err!=nil {return userlib.UUIDNew(), err}

		toRecipientInfo = ToRecipientInfo{FilelockEncKey: filelockEncKey, FilelockMacKey: filelockMacKey, 
			FilelockCipherUUID: filelockCipherUUID}
	}
	// 5. Store invitationCipher in userdata's filemap
	if userdata.FilenameToRecipient[filename] == nil {
		userdata.FilenameToRecipient[filename] = make(map[string]ToRecipientInfo)
	}
	userdata.FilenameToRecipient[filename][recipientName] = toRecipientInfo

	// 6. Encrypt updated User to User Cipher. Store back to Datastore
	userdata.StoreUser()

	// 7. Return invitation pointer
	return newInvitationCipherUUID, nil
}

/* ============================ Accept Invitation ============================== */
func (userdata *User) AcceptInvitation(senderName string, invitationPtr userlib.UUID, filename string) error {
	// 1. Reload User struct to update changes form other sessions
	var user_error error
	userdata, user_error = GetUser((*userdata).Username, (*userdata).Password)
	if user_error != nil {return user_error}

	_, ok := (*userdata).FilenameToSender[filename]
	if ok {return errors.New(strings.ToTitle("File already exists."))}

	// 2. Verify integrity of invitation cipher
	filelockEncKey, filelockMacKey, filelockCipherUUID, err := userdata.VerifyAndDecryptInvitation(invitationPtr, senderName) 
	if err!=nil {return err}

	// 3. Check deciphered invitation is correct
	filelock, err := VerifyAndDecryptFilelock(filelockCipherUUID, filelockEncKey, filelockMacKey)
	if err!=nil {return err}

	// 4. Check invitation is still valid to decrypt file (no revoke called yet)
	_, err = VerifyAndDecryptFile(filelock)
	if err!=nil {return err}

	// 5. Store invitationPtr in user's filemap
	fromSenderInfo := FromSenderInfo{SenderName: senderName, UUID: invitationPtr}
	(*userdata).FilenameToSender[filename] = fromSenderInfo

	// 6. Store FromSenderInfo with UUID(H(username||filename))
	userdata.CreateAndStoreFromSenderInfo(invitationPtr, filename, senderName)

	// 7. Encrypt updated User to User Cipher. Store back to Datastore
	userdata.StoreUser()

	return nil
}

/* ======================== Revoke Access =================================== */
func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {
	// 1 Reload userdata to update change from other sessions
	var user_error error
	userdata, user_error = GetUser((*userdata).Username, (*userdata).Password)
	if user_error != nil {return user_error}

	// 2 Check if filename is in user's filemap, if not, error
    infoFromSender, ok := userdata.FilenameToSender[filename]
	if !ok {
		return errors.New(strings.ToTitle("File not found!"))
	}

	// 3 Check invitation integrity and decrypt invitation to get filelock info.
	invitationCipherUUID := infoFromSender.UUID
	filelockEncKey, filelockMacKey, filelockCipherUUID, err := userdata.VerifyAndDecryptInvitation(invitationCipherUUID, infoFromSender.SenderName)
	if err!=nil {return err}

	// 4 Check filelock integrity and decrypt filelock to get fileEncKey & fileMacKey.
	filelock, err := VerifyAndDecryptFilelock(filelockCipherUUID, filelockEncKey, filelockMacKey)
	if err!=nil {return err}
	fileEncKey := (*filelock).FileEncKey
	fileMacKey := (*filelock).FileMacKey

	// 5 Check file integrity and decrypt file.
	file, err := VerifyAndDecryptFile(filelock)
	if err!=nil {return err}

	// 6 Generate new fileEncKey & fileMacKey
	newFileEncKey := userlib.RandomBytes(16)
	newFileMacKey := userlib.RandomBytes(16)

	// 7 Read & re-encrypt all file blocks with new fileEncKey & fileMacKey, and store them back to dataStore with new UUID.
	// Update index map in file acccordingly. Delete original file in dataStore.
	for i:=0; i<(*file).NumOfBlock; i++ {
		// 7.1 Check file block integrity and decrypt file block.
		index_byte := IntToByte(i)
		fileBlockCipherUUID, _ := userlib.UUIDFromBytes(userlib.Hash(append(fileEncKey, index_byte...)))
		fileBlock, err := VerifyAndDecryptFileBlock(fileBlockCipherUUID, fileEncKey, fileMacKey)
		if err!=nil {return err}

		// 7.2 Re-encrypt and store back to dataStore
		content := (*fileBlock).Content
		CreateAndStoreFileBlock(content, newFileEncKey, newFileMacKey, i)

		// 7.3 Delete original file block stored in dataStore
		userlib.DatastoreDelete(fileBlockCipherUUID)
	}

	// 8 Re-encrypt file with new fileEncKey & fileMacKey and store it back to dataStore with new UUID.
	// Delete original file in dataStore.
	newFileCipherUUID := userlib.UUIDNew()
	StoreFile(newFileCipherUUID, file, newFileEncKey, newFileMacKey)
	userlib.DatastoreDelete((*filelock).FileCipherUUID) //prevent revoke attack

	// 9 For all file recipients except the revoked user, 
	// update fileEncKey & fileMacKey & file cipher UUID in corresponding filelock
	found := false
	for recipient, toRecipientInfo := range (*userdata).FilenameToRecipient[filename]{
		if recipient == recipientUsername {
			found = true
			delete((*userdata).FilenameToRecipient, filename)
			continue
		}
		updatedFilelock := Filelock{FileEncKey: newFileEncKey, FileMacKey: newFileMacKey, 
			FileCipherUUID: newFileCipherUUID}
		recipientFilelockEncKey := toRecipientInfo.FilelockEncKey
		recipientFilelockMacKey := toRecipientInfo.FilelockMacKey
		recipientFilelockCipherUUID := toRecipientInfo.FilelockCipherUUID
		StoreFilelock(recipientFilelockCipherUUID, &updatedFilelock, recipientFilelockEncKey, recipientFilelockMacKey)
	}
	if !found {return errors.New(strings.ToTitle("File wasn't shared with the user!"))}
	
	// 10 Encrypt updated User to User Cipher. Store back to Datastore
	userdata.StoreUser()

	return nil
}
