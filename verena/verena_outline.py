# Verena client
def init_user(username, password):
    # send request to server to authenticate
    pass

def store_file(filename, content) -> bool:
    # send request to server to store file with content
    # verify root hash from server using proof (also check against root hash in hash server)
    pass

def load_file(filename):
    # send request to server to get filename
    # verify root hash from server using proof (also check against root hash in hash server)
    # return file
    pass


# Verena server (since the merkle tree code is in rust, this should probably be in rust)
def init_user(username, password):
    # authenticate request for username with password
    pass

def store_file(user, filename, content):
    #  store filename with content for user
    # 1. call cool-bean#LoadFile(user, filename)
    # 2. if exists, get root hash
    # 3. else, call cool-bean#StoreFile(user, filename, content)
    # compute hash for file
    # write to hash server
    # output blob:
    #   root hash
    #   inclusion proof (intermediate hashes to show that the tree is valid)
    pass

def get_file(user, filename):
    # call cool-bean#LoadFile(user, filename)
    # output blob:
    #   file
    #   root hash
    #   inclusion proof
    pass
