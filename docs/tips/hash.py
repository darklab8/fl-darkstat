# import the hash algorithm
from passlib.hash import pbkdf2_sha256

# generate new salt, and hash a password
hash_ = pbkdf2_sha256.hash("toomanysecrets")
print(hash_)
hash_2 = pbkdf2_sha256.hash("toomanysecrets")
print(hash_2)

# verifying the password
print(pbkdf2_sha256.verify("toomanysecrets", hash_))
print(pbkdf2_sha256.verify("toomanysecrets", hash_2))
print(pbkdf2_sha256.verify("joshua", hash_))
