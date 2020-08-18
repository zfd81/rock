# dasctl

## Key-value commands

### USER \<subcommand\>

USER provides commands for managing users of das.

### USER ADD \<user name or user:password\> [options]

`user add` creates a user.

RPC: UserAdd

#### Options

- interactive -- Read password from stdin instead of interactive terminal

#### Output

`User <user name> created`.

#### Examples

```bash
./dasctl --user=root:123 user add zfd
# Password of myuser: #type password for zfd
# Type password of zfd again for confirmation:#re-type password for zfd
# User zfd created
```

### USER GET \<user name\> [options]

`user get` lists detailed user information.

RPC: UserGet

#### Options

- detail -- Show permissions of roles granted to the user

#### Output

Detailed user information.

#### Examples

```bash
./dasctl --user=root:123 user get myuser
# User: myuser
# Roles:
```

### USER DELETE \<user name\>

`user delete` deletes a user.

RPC: UserDelete

#### Output

`User <user name> deleted`.

#### Examples

```bash
./dasctl --user=root:123 user delete myuser
# User myuser deleted
```

### USER LIST

`user list` lists detailed user information.

RPC: UserList

#### Output

- List of users, one per line.

#### Examples

```bash
./dasctl --user=root:123 user list
# user1
# user2
# myuser
```

### USER PASSWD \<user name\> [options]

`user passwd` changes a user's password.

RPC: UserChangePassword

#### Options

- interactive -- if true, read password in interactive terminal

#### Output

`Password updated`.

#### Examples

```bash
./dasctl --user=root:123 user passwd myuser
# Password of myuser: #type new password for my user
# Type password of myuser again for confirmation: #re-type the new password for my user
# Password updated
```

### USER GRANT-ROLE \<user name\> \<role name\>

`user grant-role` grants a role to a user

RPC: UserGrantRole

#### Output

`Role <role name> is granted to user <user name>`.

#### Examples

```bash
./dasctl --user=root:123 user grant-role userA roleA
# Role roleA is granted to user userA
```

### USER REVOKE-ROLE \<user name\> \<role name\>

`user revoke-role` revokes a role from a user

RPC: UserRevokeRole

#### Output

`Role <role name> is revoked from user <user name>`.

#### Examples

```bash
./dasctl --user=root:123 user revoke-role userA roleA
# Role roleA is revoked from user userA
```

## Utility commands
